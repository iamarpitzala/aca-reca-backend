-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Normalized Entry Tables for Net and Gross Methods
-- =====================================================

-- 1. Entry Header Table (replaces common fields from tbl_custom_form_entry)
CREATE TABLE IF NOT EXISTS tbl_entry_header (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    form_name VARCHAR(255) NOT NULL,
    form_type VARCHAR(50) NOT NULL CHECK (form_type IN ('income', 'expense', 'both')),
    calculation_method VARCHAR(50) NOT NULL CHECK (calculation_method IN ('net', 'gross')),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    quarter_id UUID REFERENCES tbl_quarter(id) ON DELETE SET NULL,
    entry_date DATE NOT NULL,
    description TEXT,
    remarks TEXT,
    payment_responsibility VARCHAR(50) CHECK (payment_responsibility IN ('owner', 'clinic')),
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    -- Link to original entry if migrating from JSONB structure
    original_entry_id UUID REFERENCES tbl_custom_form_entry(id) ON DELETE SET NULL
);

CREATE INDEX idx_entry_header_form_id ON tbl_entry_header(form_id);
CREATE INDEX idx_entry_header_clinic_id ON tbl_entry_header(clinic_id);
CREATE INDEX idx_entry_header_quarter_id ON tbl_entry_header(quarter_id);
CREATE INDEX idx_entry_header_entry_date ON tbl_entry_header(entry_date);
CREATE INDEX idx_entry_header_calculation_method ON tbl_entry_header(calculation_method);
CREATE INDEX idx_entry_header_deleted_at ON tbl_entry_header(deleted_at);
CREATE INDEX idx_entry_header_original_entry_id ON tbl_entry_header(original_entry_id);

-- 2. Entry Field Values Table (normalizes JSONB values array)
CREATE TABLE IF NOT EXISTS tbl_entry_field_value (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    field_id VARCHAR(255) NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    value NUMERIC(14,2), -- For number/currency fields
    text_value TEXT, -- For text fields (if needed in future)
    boolean_value BOOLEAN, -- For boolean fields (if needed in future)
    manual_gst_amount NUMERIC(14,2), -- When GST type is 'manual'
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_field_value_entry_id ON tbl_entry_field_value(entry_id);
CREATE INDEX idx_entry_field_value_field_id ON tbl_entry_field_value(field_id);

-- 3. Entry Field Calculations Table (normalizes fieldTotals from calculations JSONB)
CREATE TABLE IF NOT EXISTS tbl_entry_field_calculation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    field_id VARCHAR(255) NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    gst_type VARCHAR(20) NOT NULL CHECK (gst_type IN ('inclusive', 'exclusive', 'manual')),
    section VARCHAR(50), -- 'income', 'expense', 'additional_reduction'
    payment_responsibility VARCHAR(50) CHECK (payment_responsibility IN ('owner', 'clinic')),
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_field_calc_entry_id ON tbl_entry_field_calculation(entry_id);
CREATE INDEX idx_entry_field_calc_field_id ON tbl_entry_field_calculation(field_id);
CREATE INDEX idx_entry_field_calc_section ON tbl_entry_field_calculation(section);
CREATE INDEX idx_entry_field_calc_payment_resp ON tbl_entry_field_calculation(payment_responsibility);

-- 4. Entry Summary Calculations Table (normalizes summary totals from calculations JSONB)
CREATE TABLE IF NOT EXISTS tbl_entry_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    -- Common totals
    total_base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_payable NUMERIC(14,2) NOT NULL DEFAULT 0, -- For expense forms
    net_receivable NUMERIC(14,2) NOT NULL DEFAULT 0, -- For income forms
    net_fee NUMERIC(14,2), -- For income/both forms
    -- BAS mapping
    bas_gst_on_sales_1a NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_gst_credit_1b NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_total_sales_g1 NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_expenses_g11 NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_summary_entry_id ON tbl_entry_summary(entry_id);

-- 5. Net Method Entry Details Table (for independent contractors)
CREATE TABLE IF NOT EXISTS tbl_entry_net_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    commission_percent NUMERIC(5,2) NOT NULL,
    commission NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_commission NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_payment_received NUMERIC(14,2) NOT NULL DEFAULT 0,
    -- Super holding fields (when enabled)
    super_holding_enabled BOOLEAN NOT NULL DEFAULT false,
    super_component_percent NUMERIC(5,2), -- Default 12%
    commission_component NUMERIC(14,2), -- When super holding enabled
    super_component NUMERIC(14,2), -- When super holding enabled
    total_for_reconciliation NUMERIC(14,2), -- When super holding enabled
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_net_details_entry_id ON tbl_entry_net_details(entry_id);

-- 6. Gross Method Entry Details Table (for service/facility fees)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    service_facility_fee_percent NUMERIC(5,2) NOT NULL,
    service_fee_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    subtotal_after_deductions NUMERIC(14,2),
    remitted_amount NUMERIC(14,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_details_entry_id ON tbl_entry_gross_details(entry_id);

-- 7. Gross Method Reductions Table (clinic-paid expenses - GST portion)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    field_calculation_id UUID REFERENCES tbl_entry_field_calculation(id) ON DELETE CASCADE,
    field_id VARCHAR(255) NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0, -- Only GST portion is shown in reductions
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_reduction_entry_id ON tbl_entry_gross_reduction(entry_id);
CREATE INDEX idx_entry_gross_reduction_field_calc_id ON tbl_entry_gross_reduction(field_calculation_id);

-- 8. Gross Method Reimbursements Table (owner-paid expenses)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reimbursement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    field_calculation_id UUID REFERENCES tbl_entry_field_calculation(id) ON DELETE CASCADE,
    field_id VARCHAR(255) NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_reimbursement_entry_id ON tbl_entry_gross_reimbursement(entry_id);
CREATE INDEX idx_entry_gross_reimbursement_field_calc_id ON tbl_entry_gross_reimbursement(field_calculation_id);

-- 9. Gross Method Additional Reductions Table (Merchant Fee, Bank Fee, etc.)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_additional_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    field_calculation_id UUID REFERENCES tbl_entry_field_calculation(id) ON DELETE CASCADE,
    field_id VARCHAR(255) NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_addl_reduction_entry_id ON tbl_entry_gross_additional_reduction(entry_id);
CREATE INDEX idx_entry_gross_addl_reduction_field_calc_id ON tbl_entry_gross_additional_reduction(field_calculation_id);

-- 10. Gross Method Reductions Summary Table
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reductions_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    total_reductions NUMERIC(14,2) NOT NULL DEFAULT 0, -- Effective reductions (GST or outwork charge)
    total_reduction_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_expense_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_reimbursements NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_reductions_summary_entry_id ON tbl_entry_gross_reductions_summary(entry_id);

-- 11. Gross Method Outwork Charge Table (consolidated expense charge)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_outwork (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    outwork_enabled BOOLEAN NOT NULL DEFAULT false,
    outwork_rate_percent NUMERIC(5,2),
    outwork_charge_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    outwork_charge_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
    outwork_charge_total NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_gross_outwork_entry_id ON tbl_entry_gross_outwork(entry_id);

-- 12. Entry Deductions Table (stores deduction parameters used at entry time)
CREATE TABLE IF NOT EXISTS tbl_entry_deductions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL UNIQUE REFERENCES tbl_entry_header(id) ON DELETE CASCADE,
    -- Service fee parameters (for gross method)
    service_facility_fee_percent NUMERIC(5,2),
    service_fee_override NUMERIC(14,2),
    -- Commission parameters (for net method)
    commission_percent NUMERIC(5,2),
    super_holding_enabled BOOLEAN DEFAULT false,
    super_component_percent NUMERIC(5,2), -- Default 12%
    -- Outwork parameters (for gross method)
    outwork_enabled BOOLEAN DEFAULT false,
    outwork_rate_percent NUMERIC(5,2),
    -- Entry-level payment responsibility override
    entry_payment_responsibility VARCHAR(50) CHECK (entry_payment_responsibility IN ('owner', 'clinic')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_entry_deductions_entry_id ON tbl_entry_deductions(entry_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS tbl_entry_deductions;
DROP TABLE IF EXISTS tbl_entry_gross_outwork;
DROP TABLE IF EXISTS tbl_entry_gross_reductions_summary;
DROP TABLE IF EXISTS tbl_entry_gross_additional_reduction;
DROP TABLE IF EXISTS tbl_entry_gross_reimbursement;
DROP TABLE IF EXISTS tbl_entry_gross_reduction;
DROP TABLE IF EXISTS tbl_entry_gross_details;
DROP TABLE IF EXISTS tbl_entry_net_details;
DROP TABLE IF EXISTS tbl_entry_summary;
DROP TABLE IF EXISTS tbl_entry_field_calculation;
DROP TABLE IF EXISTS tbl_entry_field_value;
DROP TABLE IF EXISTS tbl_entry_header;

-- +goose StatementEnd
