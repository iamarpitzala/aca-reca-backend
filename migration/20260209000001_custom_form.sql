-- +goose Up
-- +goose StatementBegin
-- Custom form definitions (builder schema)
CREATE TABLE IF NOT EXISTS tbl_custom_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    calculation_method VARCHAR(50) NOT NULL CHECK (calculation_method IN ('net', 'gross')),
    form_type VARCHAR(50) NOT NULL CHECK (form_type IN ('income', 'expense', 'both')),
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    fields JSONB NOT NULL DEFAULT '[]',
    default_payment_responsibility VARCHAR(50) CHECK (default_payment_responsibility IN ('owner', 'clinic')),
    service_facility_fee_percent NUMERIC(5,2),
    version INTEGER NOT NULL DEFAULT 1,
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_custom_form_clinic_id ON tbl_custom_form(clinic_id);
CREATE INDEX idx_custom_form_status ON tbl_custom_form(status);
CREATE INDEX idx_custom_form_deleted_at ON tbl_custom_form(deleted_at);

-- Note: Custom form entries are now stored in normalized tables (tbl_entry_header, etc.)
-- See migration 20260214000000_normalized_entry_tables.sql
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form;
-- +goose StatementEnd
