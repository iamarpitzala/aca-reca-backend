-- +goose Up
-- +goose StatementBegin

-- Entry Summary Table
CREATE TABLE IF NOT EXISTS tbl_entry_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    total_base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_payable NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_receivable NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_fee NUMERIC(14,2),
    bas_gst_on_sales_1a NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_gst_credit_1b NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_total_sales_g1 NUMERIC(14,2) NOT NULL DEFAULT 0,
    bas_expenses_g11 NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Net Details Table
CREATE TABLE IF NOT EXISTS tbl_entry_net_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
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

-- Gross Details Table (service/facility fees)
CREATE TABLE IF NOT EXISTS tbl_entry_gross_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    service_facility_fee_percent NUMERIC(5,2) NOT NULL,
    service_fee_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    subtotal_after_deductions NUMERIC(14,2),
    remitted_amount NUMERIC(14,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Reduction Table
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Reimbursement Table
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reimbursement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Additional Reduction Table
CREATE TABLE IF NOT EXISTS tbl_entry_gross_additional_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Reductions Summary Table
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reductions_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    total_reductions NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_reduction_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_expense_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_reimbursements NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction_gst NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_additional_reduction_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- Entry Deductions Table (stores deduction parameters used at entry time)
CREATE TABLE IF NOT EXISTS tbl_entry_deductions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    -- Service fee parameters (for gross method)
    service_facility_fee_percent NUMERIC(5,2),
    service_fee_override NUMERIC(14,2),
    -- Commission parameters (for net method)
    commission_percent NUMERIC(5,2),
    super_holding_enabled BOOLEAN DEFAULT false,
    super_component_percent NUMERIC(5,2), -- Default 12%
    super_component NUMERIC(14,2), -- When super holding enabled
    total_for_reconciliation NUMERIC(14,2), -- When super holding enabled
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS tbl_entry_deductions;
DROP TABLE IF EXISTS tbl_entry_gross_reductions_summary;
DROP TABLE IF EXISTS tbl_entry_gross_additional_reduction;
DROP TABLE IF EXISTS tbl_entry_gross_reimbursement;
DROP TABLE IF EXISTS tbl_entry_gross_reduction;
DROP TABLE IF EXISTS tbl_entry_gross_details;
DROP TABLE IF EXISTS tbl_entry_net_details;
DROP TABLE IF EXISTS tbl_entry_summary;

-- +goose StatementEnd
