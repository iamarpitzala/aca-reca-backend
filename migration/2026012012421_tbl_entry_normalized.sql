-- +goose Up
-- +goose StatementBegin

-- Entry Summary
CREATE TABLE IF NOT EXISTS tbl_entry_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    total_base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_payable NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_receivable NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_fee NUMERIC(14,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Net Details
CREATE TABLE IF NOT EXISTS tbl_entry_net_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
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

-- Gross Details
CREATE TABLE IF NOT EXISTS tbl_entry_gross_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entry_id UUID NOT NULL UNIQUE REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    service_facility_fee_percent NUMERIC(5,2) NOT NULL,
    service_fee_base NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_on_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_service_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Reduction
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reduction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gross_details_id UUID NOT NULL REFERENCES tbl_entry_gross_details(id) ON DELETE CASCADE,
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Gross Reimbursement
CREATE TABLE IF NOT EXISTS tbl_entry_gross_reimbursement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gross_details_id UUID NOT NULL REFERENCES tbl_entry_gross_details(id) ON DELETE CASCADE,
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    base_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_entry_gross_reduction;
DROP TABLE IF EXISTS tbl_entry_gross_reimbursement;
DROP TABLE IF EXISTS tbl_entry_gross_details;
DROP TABLE IF EXISTS tbl_entry_summary;
DROP TABLE IF EXISTS tbl_entry_net_details;

-- +goose StatementEnd
