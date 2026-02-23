# Analysis: `tbl_entry_net_details` Table

## Table Schema

```sql
CREATE TABLE IF NOT EXISTS tbl_entry_net_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id),
    commission_percent NUMERIC(5,2) NOT NULL,
    commission NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_commission NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_payment_received NUMERIC(14,2) NOT NULL DEFAULT 0,
    super_holding_enabled BOOLEAN NOT NULL DEFAULT false,
    super_component_percent NUMERIC(5,2),
    commission_component NUMERIC(14,2),
    super_component NUMERIC(14,2),
    total_for_reconciliation NUMERIC(14,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Purpose

The `tbl_entry_net_details` table stores **NET calculation method** specific details for form entries. This table is used when a form's `calculation_method` is set to **"NET"**.

## Field Analysis

### Primary Key & Relationships
- **`id`**: Primary key (UUID)
- **`source_entry_id`**: Foreign key to `tbl_custom_form_entry(id)` - ONE-TO-ONE relationship (UNIQUE constraint)
  - Links to the main entry record
  - CASCADE DELETE ensures net details are deleted when entry is deleted

### Commission Fields
- **`commission_percent`**: Commission percentage (NUMERIC 5,2) - **REQUIRED**
  - Example: 30.00 = 30%
  - Used to calculate commission from total payment
  
- **`commission`**: Calculated commission amount (NUMERIC 14,2)
  - Formula: `commission = total_payment_received * (commission_percent / 100)`
  
- **`gst_on_commission`**: GST amount on commission (NUMERIC 14,2)
  - GST calculated on the commission amount
  - Typically: `gst_on_commission = commission * gst_rate / (100 + gst_rate)` for inclusive GST

### Payment Fields
- **`total_payment_received`**: Total payment amount received (NUMERIC 14,2)
  - Base amount before commission/super deductions
  - This is the gross payment received

### Superannuation (Super) Fields
- **`super_holding_enabled`**: Boolean flag indicating if super holding is enabled
  - Default: `false`
  - When `true`, superannuation component is calculated and held
  
- **`super_component_percent`**: Superannuation percentage (NUMERIC 5,2)
  - Example: 10.50 = 10.5%
  - Only used when `super_holding_enabled = true`
  
- **`commission_component`**: Commission portion of payment (NUMERIC 14,2)
  - Part of total payment that goes to commission
  
- **`super_component`**: Superannuation component amount (NUMERIC 14,2)
  - Calculated when super holding is enabled
  - Formula: `super_component = commission_component * (super_component_percent / 100)`

### Reconciliation Field
- **`total_for_reconciliation`**: Final amount for reconciliation (NUMERIC 14,2)
  - Net amount after all deductions
  - Formula: `total_for_reconciliation = commission_component - super_component`

### Timestamps
- **`created_at`**: Record creation timestamp
- **`updated_at`**: Last update timestamp

## Relationship to Calculation Method

This table is **ONLY populated** when:
1. Form's `calculation_method = 'NET'` (from `tbl_custom_form` or `tbl_custom_form_calculation`)
2. Entry is created with NET calculation method
3. Net-specific calculations are performed

## Comparison: NET vs GROSS

### NET Method (`tbl_entry_net_details`)
- Focuses on **commission-based** calculations
- Calculates commission from total payment received
- Handles superannuation holding
- Used for **dentist commission** scenarios

### GROSS Method (`tbl_entry_gross_details`)
- Focuses on **service facility fee** calculations
- Calculates service fee from net amount
- Handles reductions and reimbursements
- Used for **clinic fee** scenarios

## Current Implementation Status

### ✅ What Exists
1. **Database Schema**: Table is created via migration `2026012012421_tbl_entry_normalized.sql`
2. **Table Structure**: Properly defined with all necessary fields
3. **Relationships**: Foreign key to `tbl_custom_form_entry` with CASCADE DELETE

### ❌ What's Missing
1. **Domain Model**: No `EntryNetDetails` struct in `internal/domain/entry_normalized.go`
   - `EntryGrossDetails` exists, but `EntryNetDetails` is missing
   
2. **Repository Interface**: No repository methods to:
   - Create net details
   - Get net details by entry ID
   - Update net details
   
3. **Calculation Logic**: No NET calculation implementation
   - `RunGrossCalculation()` exists in `internal/calculation/gross.go`
   - No equivalent `RunNetCalculation()` function
   - Current `RunEntryCalculation()` in `custom_form.go` only handles GROSS or basic NET (without net details)

4. **Use Case Integration**: Net details are not being:
   - Calculated during entry creation
   - Stored in database
   - Retrieved when fetching entries

## Recommended Implementation Steps

### Step 1: Create Domain Model
Add to `internal/domain/entry_normalized.go`:
```go
// EntryNetDetails represents tbl_entry_net_details
type EntryNetDetails struct {
    ID                      uuid.UUID `db:"id"`
    EntryID                 uuid.UUID `db:"source_entry_id"`
    CommissionPercent       float64   `db:"commission_percent"`
    Commission              float64   `db:"commission"`
    GSTOnCommission         float64   `db:"gst_on_commission"`
    TotalPaymentReceived     float64   `db:"total_payment_received"`
    SuperHoldingEnabled      bool      `db:"super_holding_enabled"`
    SuperComponentPercent    *float64  `db:"super_component_percent"`
    CommissionComponent      *float64  `db:"commission_component"`
    SuperComponent           *float64  `db:"super_component"`
    TotalForReconciliation   *float64  `db:"total_for_reconciliation"`
    CreatedAt                time.Time `db:"created_at"`
    UpdatedAt                time.Time `db:"updated_at"`
}
```

### Step 2: Create Repository
- Add `EntryNetDetailsRepository` interface in `internal/application/port/`
- Implement in `internal/adapter/postgres/`
- Methods: `Create()`, `GetByEntryID()`, `Update()`, `Delete()`

### Step 3: Implement NET Calculation Logic
- Create `internal/calculation/net.go` with `RunNetCalculation()` function
- Calculate commission, GST on commission, super components
- Return structured output similar to gross calculation

### Step 4: Integrate with Entry Creation
- Update `CreateEntry()` in `field_entry.go` to:
  - Detect NET calculation method
  - Call NET calculation engine
  - Save net details to `tbl_entry_net_details`

### Step 5: Update Entry Retrieval
- Modify `GetByID()`, `GetByFormID()` to include net details
- Add net details to response structure

## Data Flow (When Implemented)

```
Entry Creation (NET Method)
    ↓
Calculate Totals (NET)
    ↓
Calculate Commission & Super
    ↓
Save to tbl_entry_net_details
    ↓
Return EntryResponse with Net Details
```

## Key Formulas (NET Method)

1. **Commission Calculation**:
   ```
   commission = total_payment_received × (commission_percent / 100)
   ```

2. **GST on Commission** (if GST inclusive):
   ```
   gst_on_commission = commission × (gst_rate / (100 + gst_rate))
   ```
   (if GST exclusive):
   ```
   gst_on_commission = commission × (gst_rate / 100)
   ```

3. **Super Component** (if enabled):
   ```
   super_component = commission_component × (super_component_percent / 100)
   ```

4. **Total for Reconciliation**:
   ```
   total_for_reconciliation = commission_component - super_component
   ```

## Notes

- The table uses `source_entry_id` (not `entry_id`) to match naming convention with gross details
- `commission_percent` is NOT NULL, meaning it must always be provided
- Super fields are nullable, only populated when super holding is enabled
- This table is part of the "normalized entry" structure for better data organization
