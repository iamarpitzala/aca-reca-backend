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

-- Custom form entries (submitted data)
CREATE TABLE IF NOT EXISTS tbl_custom_form_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    form_name VARCHAR(255) NOT NULL,
    form_type VARCHAR(50) NOT NULL CHECK (form_type IN ('income', 'expense', 'both')),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    quarter_id UUID REFERENCES tbl_quarter(id) ON DELETE SET NULL,
    values JSONB NOT NULL DEFAULT '[]',
    calculations JSONB NOT NULL DEFAULT '{}',
    entry_date DATE NOT NULL,
    description TEXT,
    remarks TEXT,
    payment_responsibility VARCHAR(50) CHECK (payment_responsibility IN ('owner', 'clinic')),
    deductions JSONB,
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_custom_form_entry_form_id ON tbl_custom_form_entry(form_id);
CREATE INDEX idx_custom_form_entry_clinic_id ON tbl_custom_form_entry(clinic_id);
CREATE INDEX idx_custom_form_entry_quarter_id ON tbl_custom_form_entry(quarter_id);
CREATE INDEX idx_custom_form_entry_entry_date ON tbl_custom_form_entry(entry_date);
CREATE INDEX idx_custom_form_entry_deleted_at ON tbl_custom_form_entry(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_entry;
DROP TABLE IF EXISTS tbl_custom_form;
-- +goose StatementEnd
