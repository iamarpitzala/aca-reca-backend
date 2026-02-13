-- +goose Up
-- +goose StatementBegin
-- Custom form definitions (builder schema)
CREATE TABLE tbl_custom_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    form_type VARCHAR(50) NOT NULL CHECK (form_type IN ('INCOME', 'EXPENSE', 'BOTH')),
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED')),
    calculation_method VARCHAR(50) NOT NULL CHECK (calculation_method IN ('NET', 'GROSS')),
    default_payment_responsibility VARCHAR(50) NOT NULL CHECK (default_payment_responsibility IN ('OWNER', 'CLINIC')),
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE tbl_custom_form_version (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    UNIQUE (form_id, version)
);


CREATE TABLE tbl_custom_form_field (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_version_id UUID NOT NULL REFERENCES tbl_custom_form_version(id) ON DELETE CASCADE,
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    field_key VARCHAR(100) NOT NULL,
    label VARCHAR(255) NOT NULL,
    field_type VARCHAR(50) NOT NULL,
    section VARCHAR(50) NOT NULL CHECK (section IN ('INCOME', 'EXPENSE', 'REDUCTION')),
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    coa_id UUID NOT NULL REFERENCES tbl_account(id) ON DELETE RESTRICT,
    placeholder VARCHAR(255) NOT NULL,
    min_value NUMERIC(14,2),
    max_value NUMERIC(14,2),
    field_order INTEGER NOT NULL,
    gst_config BOOLEAN NOT NULL DEFAULT FALSE,
    gst_rate NUMERIC(5,2),
    gst_type VARCHAR(50) NOT NULL CHECK (gst_type IN ('INCLUSIVE', 'EXCLUSIVE', 'MANUAL')),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (form_version_id, field_key)
);


CREATE TABLE tbl_custom_form_calculation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_version_id UUID NOT NULL REFERENCES tbl_custom_form_version(id) ON DELETE CASCADE,
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    calculation_method VARCHAR(50) NOT NULL
        CHECK (calculation_method IN ('NET', 'GROSS')),
    default_payment_responsibility VARCHAR(50)
        CHECK (default_payment_responsibility IN ('OWNER', 'CLINIC')),
    service_facility_fee_percent NUMERIC(5,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE tbl_custom_form_publish (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_version_id UUID NOT NULL REFERENCES tbl_custom_form_version(id),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    published_by UUID NOT NULL REFERENCES tbl_user(id),
    published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Note: Custom form entries are now stored in normalized tables (tbl_entry_header, etc.)
-- See migration 20260214000000_normalized_entry_tables.sql
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form;
-- +goose StatementEnd
