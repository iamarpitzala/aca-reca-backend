# Analysis: GROSS Tables and Entry Creation Flow

## 1. GROSS Table Structure Analysis

### `tbl_entry_gross_details`
**Purpose**: Stores GROSS method calculation results for each entry

**Current Schema** (from migration `2026012012421_tbl_entry_normalized.sql`):
```sql
CREATE TABLE IF NOT EXISTS tbl_entry_gross_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id),
    service_facility_fee_percent NUMERIC(5,2) NOT NULL,
    service_fee_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Missing Fields** (need to be added):
- `income NUMERIC(14,2)` - Total income (including GST)
- `gst_on_income NUMERIC(14,2)` - GST amount on income
- `income_excl_gst NUMERIC(14,2)` - Income excluding GST
- `net_amount NUMERIC(14,2)` - Net amount (income_excl_gst - total_net_expenses)

### `tbl_entry_gross_reduction`
**Purpose**: Stores reduction fields (expenses) associated with GROSS details

**Schema**:
```sql
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gross_details_id UUID NOT NULL REFERENCES tbl_entry_gross_details(id),
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id),
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id),
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Fields**:
- `gross_details_id`: Links to parent GROSS details record
- `source_entry_id`: Links to the entry
- `tbl_custom_form_field_id`: Links to the form field
- `base_amount`: Base amount (excluding GST)
- `gst_amount`: GST amount
- `total_amount`: Total amount (base + GST)

### `tbl_entry_gross_reimbursement`
**Purpose**: Stores reimbursement fields (remitted costs) associated with GROSS details

**Schema**:
```sql
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reimbursement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gross_details_id UUID NOT NULL REFERENCES tbl_entry_gross_details(id),
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id),
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id),
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 2. Current Implementation Status

### ✅ Implemented:
1. **Structured GROSS Calculation** (`internal/calculation/net.go`):
   - `GrossCalculationInput` struct
   - `GrossCalculationOutput` struct
   - `RunGrossCalculationStructured()` function
   - `ParseGrossDeductions()` function

2. **Repository Interfaces** (`internal/application/port/field_entry.go`):
   - `EntryGrossDetailsRepository`
   - `EntryGrossReductionRepository`
   - `EntryGrossReimbursementRepository`

3. **Repository Implementations** (`internal/adapter/postgres/`):
   - `entry_gross_details.go` - CRUD operations
   - `entry_gross_reduction.go` - Batch create, get, delete
   - `entry_gross_reimbursement.go` - Batch create, get, delete

4. **Calculation Function** (`internal/application/usecase/field_entry.go`):
   - `calculateAndStoreGrossDetails()` function exists

### ❌ Missing/Issues:

1. **GROSS Method Check Missing in CreateEntry**:
   - Line 224-230: Only checks for NET method
   - **Missing**: Check for GROSS method and call `calculateAndStoreGrossDetails()`

2. **Domain Model Missing Income Fields**:
   - `EntryGrossDetails` struct doesn't have:
     - `Income`
     - `GstOnIncome`
     - `IncomeExclGst`
     - `NetAmount`

3. **Database Migration Missing**:
   - Need migration to add income fields to `tbl_entry_gross_details`

4. **Repository Query Missing Fields**:
   - `entry_gross_details.go` needs to include income fields in INSERT/UPDATE queries

## 3. Entry Creation Flow (Current vs Expected)

### Current Flow:
```
CreateEntry()
  ↓
Calculate Entry Totals (RunEntryCalculation)
  ↓
Store Entry
  ↓
Check if NET method → calculateAndStoreNetDetails()
  ↓
❌ MISSING: Check if GROSS method → calculateAndStoreGrossDetails()
  ↓
Return Response
```

### Expected Flow:
```
CreateEntry()
  ↓
Calculate Entry Totals (RunEntryCalculation)
  ↓
Store Entry
  ↓
Check if NET method → calculateAndStoreNetDetails()
  ↓
✅ Check if GROSS method → calculateAndStoreGrossDetails()
  ↓
Return Response
```

## 4. Required Fixes

### Fix 1: Add GROSS Method Check in CreateEntry
**Location**: `internal/application/usecase/field_entry.go` (after NET check)

```go
// If GROSS method, calculate and store gross details
if calculationMethod == "GROSS" {
    if err := s.calculateAndStoreGrossDetails(ctx, form.ClinicID, formVersionID, fieldEntries[0].ID, calculationsJSON, fieldValueResponses, req.Deductions, now); err != nil {
        // Log error but don't fail entry creation
        _ = err
    }
}
```

### Fix 2: Add Income Fields to Domain Model
**Location**: `internal/domain/entry_normalized.go`

Add to `EntryGrossDetails`:
```go
Income        float64 `db:"income"`
GstOnIncome   float64 `db:"gst_on_income"`
IncomeExclGst float64 `db:"income_excl_gst"`
NetAmount     float64 `db:"net_amount"`
```

### Fix 3: Create Migration for Income Fields
**File**: `migration/YYYYMMDDHHMMSS_add_income_fields_to_gross_details.sql`

```sql
ALTER TABLE tbl_entry_gross_details
ADD COLUMN income NUMERIC(14,2) NOT NULL DEFAULT 0,
ADD COLUMN gst_on_income NUMERIC(14,2) NOT NULL DEFAULT 0,
ADD COLUMN income_excl_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
ADD COLUMN net_amount NUMERIC(14,2) NOT NULL DEFAULT 0;
```

### Fix 4: Update Repository Queries
**Location**: `internal/adapter/postgres/entry_gross_details.go`

Update INSERT and UPDATE queries to include income fields.

## 5. Data Flow (When Fixed)

```
Entry Creation (GROSS Method)
  ↓
RunEntryCalculation (GROSS)
  ↓
calculateAndStoreGrossDetails()
  ├─ Parse deductions
  ├─ Get form calculation settings
  ├─ Get GST settings
  ├─ Parse calculation output (IncomeExclGST, NetAmount)
  ├─ Build field totals
  ├─ RunGrossCalculationStructured()
  ├─ Create EntryGrossDetails
  │   ├─ Service fee fields
  │   └─ Income fields (income, gst_on_income, income_excl_gst, net_amount)
  ├─ Save to tbl_entry_gross_details
  ├─ Map reductions (REDUCTION section fields)
  └─ Save to tbl_entry_gross_reduction
  ↓
Return EntryResponse
```

## 6. Verification Checklist

- [ ] GROSS method check added in CreateEntry
- [ ] Income fields added to domain model
- [ ] Migration created and applied
- [ ] Repository queries updated
- [ ] Test entry creation with GROSS method
- [ ] Verify data stored in tbl_entry_gross_details
- [ ] Verify reductions stored in tbl_entry_gross_reduction
- [ ] Verify income fields populated correctly
