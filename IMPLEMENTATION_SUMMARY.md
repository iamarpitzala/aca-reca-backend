# GROSS Method Implementation Summary

## ✅ Completed Implementation

### 1. GROSS Method Check Added in CreateEntry
**File**: `internal/application/usecase/field_entry.go`
- Added check for GROSS calculation method after NET check
- Calls `calculateAndStoreGrossDetails()` when method is GROSS

### 2. Structured GROSS Calculation Functions
**File**: `internal/calculation/net.go`
- ✅ `GrossCalculationInput` struct
- ✅ `GrossCalculationOutput` struct  
- ✅ `FieldTotal` struct
- ✅ `RunGrossCalculationStructured()` function
- ✅ `ParseGrossDeductions()` function

### 3. calculateAndStoreGrossDetails Function
**File**: `internal/application/usecase/field_entry.go`
- Implements structured GROSS calculation similar to NET method
- Steps:
  1. Parse deductions
  2. Get form calculation settings
  3. Get GST settings
  4. Parse calculation output (IncomeExclGST, NetAmount)
  5. Build field totals with section info
  6. Run structured GROSS calculation
  7. Create and store gross details
  8. Map and store reductions

### 4. Domain Model Updated
**File**: `internal/domain/entry_normalized.go`
- Added income fields to `EntryGrossDetails`:
  - `Income` (float64)
  - `GstOnIncome` (float64)
  - `IncomeExclGst` (float64)
  - `NetAmount` (float64)

### 5. Repository Interfaces Created
**File**: `internal/application/port/field_entry.go`
- ✅ `EntryGrossDetailsRepository`
- ✅ `EntryGrossReductionRepository`
- ✅ `EntryGrossReimbursementRepository`

### 6. Repository Implementations Created
**Files**:
- `internal/adapter/postgres/entry_gross_details.go`
- `internal/adapter/postgres/entry_gross_reduction.go`
- `internal/adapter/postgres/entry_gross_reimbursement.go`

### 7. Database Migration Created
**File**: `migration/20260216000003_add_income_fields_to_gross_details.sql`
- Adds `income`, `gst_on_income`, `income_excl_gst`, `net_amount` columns to `tbl_entry_gross_details`

### 8. Service Updated
**File**: `internal/application/usecase/field_entry.go`
- Updated `FieldEntryService` struct to include GROSS repositories
- Updated `NewFieldEntryService` constructor to accept GROSS repositories

### 9. Route Updated
**File**: `route/route.go`
- Instantiated GROSS repositories
- Injected GROSS repositories into `FieldEntryService`

## Entry Creation Flow (GROSS Method)

```
POST /api/v1/entry
  ↓
CreateEntry()
  ↓
Calculate Entry Totals (RunEntryCalculation - GROSS)
  ↓
Store Entry
  ↓
Check calculationMethod == "GROSS"
  ↓
calculateAndStoreGrossDetails()
  ├─ Parse deductions (serviceFacilityFeePercent, outworkEnabled)
  ├─ Get form calculation settings
  ├─ Get GST settings from clinic financial settings
  ├─ Parse calculation output (IncomeExclGST, NetAmount, FieldTotals)
  ├─ Get form fields (once, reused)
  ├─ Build field totals with section info
  ├─ RunGrossCalculationStructured()
  │   ├─ Calculate income (from INCOME section fields)
  │   ├─ Calculate GST on income
  │   ├─ Calculate service fee base
  │   ├─ Calculate GST on service fee
  │   └─ Calculate total service fee
  ├─ Create EntryGrossDetails
  │   ├─ Service fee fields
  │   └─ Income fields (income, gst_on_income, income_excl_gst, net_amount)
  ├─ Save to tbl_entry_gross_details
  ├─ Map reductions (REDUCTION section fields)
  └─ Save to tbl_entry_gross_reduction
  ↓
Return EntryResponse
```

## Database Tables

### tbl_entry_gross_details
Stores GROSS calculation results:
- Service fee fields: `service_facility_fee_percent`, `service_fee_base`, `gst_on_service_fee`, `total_service_fee`
- Income fields: `income`, `gst_on_income`, `income_excl_gst`, `net_amount`

### tbl_entry_gross_reduction
Stores reduction fields (expenses):
- Links to `tbl_entry_gross_details` via `gross_details_id`
- Stores `base_amount`, `gst_amount`, `total_amount` for each reduction field

### tbl_entry_gross_reimbursement
Stores reimbursement fields (remitted costs):
- Links to `tbl_entry_gross_details` via `gross_details_id`
- Stores `base_amount`, `gst_amount`, `total_amount` for each reimbursement field

## Next Steps

1. **Run Migration**: Apply the migration to add income fields to database
   ```bash
   goose -dir migration postgres "connection_string" up
   ```

2. **Test Entry Creation**: Create an entry with GROSS method and verify:
   - Data stored in `tbl_entry_gross_details`
   - Income fields populated correctly
   - Reductions stored in `tbl_entry_gross_reduction`

3. **Verify Calculations**: Check that:
   - Income is calculated from INCOME section fields
   - Service fee is calculated correctly based on outwork enabled/disabled
   - GST on income and service fee are calculated correctly

4. **Add Reimbursement Logic**: Currently only reductions are stored. Add logic to distinguish and store reimbursements if needed.

## Files Modified/Created

### Modified:
- `internal/application/usecase/field_entry.go`
- `internal/domain/entry_normalized.go`
- `internal/calculation/net.go`
- `internal/application/port/field_entry.go`
- `route/route.go`

### Created:
- `migration/20260216000003_add_income_fields_to_gross_details.sql`
- `internal/adapter/postgres/entry_gross_details.go`
- `internal/adapter/postgres/entry_gross_reduction.go`
- `internal/adapter/postgres/entry_gross_reimbursement.go`
- `ANALYSIS_GROSS_TABLES.md`
- `IMPLEMENTATION_SUMMARY.md`
