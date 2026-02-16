# Analysis: Net Amount Calculation from Field Value Responses

## Payload Structure

### EntryFieldValueResponse
```go
type EntryFieldValueResponse struct {
    FieldID         string   `json:"fieldId"`         // Field UUID as string
    FieldName       string   `json:"fieldName"`      // Field name/label
    Value           float64  `json:"value"`           // Raw value entered
    BaseAmount      *float64 `json:"baseAmount,omitempty"`  // Base amount (excl GST)
    GSTAmount       *float64 `json:"gstAmount,omitempty"`   // GST amount
    TotalAmount     *float64 `json:"totalAmount,omitempty"` // Total amount (incl GST)
    ManualGSTAmount *float64 `json:"manualGstAmount,omitempty"` // Manual GST override
}
```

### CustomFormField
```go
type CustomFormField struct {
    ID          uuid.UUID
    Section     string  // "INCOME", "EXPENSE", "REDUCTION"
    GstConfig   *GstConfig
    // ... other fields
}

type GstConfig struct {
    Enabled bool
    Rate    float64  // e.g., 10.0 for 10%
    Type    string   // "INCLUSIVE", "EXCLUSIVE", "MANUAL"
}
```

## Calculation Logic Required

### Step 1: Calculate Net Income
For each field where `Section == "INCOME"`:
1. Get field by `FieldID` from `fieldValueResponses`
2. Check `GstConfig.Type`:
   - **INCLUSIVE**: `netIncome = Value - GST`, where `GST = Value / (1 + Rate/100)`
   - **EXCLUSIVE**: `netIncome = Value` (GST is separate)
   - **MANUAL**: `netIncome = Value - ManualGSTAmount`
3. Sum all net income: `totalNetIncome += netIncome`

### Step 2: Calculate Net Expenses
For each field where `Section == "EXPENSE"`:
1. Get field by `FieldID` from `fieldValueResponses`
2. Check `GstConfig.Type`:
   - **INCLUSIVE**: `netExpense = Value - GST`, where `GST = Value / (1 + Rate/100)`
   - **EXCLUSIVE**: `netExpense = Value` (GST is separate)
   - **MANUAL**: `netExpense = Value - ManualGSTAmount`
3. Sum all net expenses: `totalNetExpenses += netExpense`

### Step 3: Calculate Net Amount
```
netAmount = totalNetIncome - totalNetExpenses
```

## Implementation Plan

1. Create function `calculateNetAmountFromFieldValues()` in `internal/calculation/net.go`
2. Function signature:
   ```go
   func CalculateNetAmountFromFieldValues(
       fieldValueResponses []domain.EntryFieldValueResponse,
       fields []domain.CustomFormField,
   ) float64
   ```
3. Logic:
   - Create field map by ID for quick lookup
   - Iterate through `fieldValueResponses`
   - For each response, find corresponding field
   - Check section (INCOME vs EXPENSE)
   - Calculate net amount based on GST type
   - Sum income and expenses separately
   - Return netIncome - netExpenses

4. Use this function in `calculateAndStoreGrossDetails()` before calling `RunGrossCalculationStructured()`
