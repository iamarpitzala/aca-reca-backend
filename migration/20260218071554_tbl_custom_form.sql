-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_tax_type (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCLUSIVE', 'EXCLUSIVE', 'MANUAL')),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

INSERT INTO tbl_tax_type (id, name, type, description) VALUES
    (1, 'Inclusive', 'INCLUSIVE', 'Inclusive tax type'),
    (2, 'Exclusive', 'EXCLUSIVE', 'Exclusive tax type'),
    (3, 'Manual', 'MANUAL', 'Manual tax type')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS tbl_section_type (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

INSERT INTO tbl_section_type (id, name, type, description) VALUES
    (1, 'INCOME', 'INCOME', 'Income section'),
    (2, 'EXPENSE', 'EXPENSE', 'Expense section')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS tbl_custom_form (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL, -- DRAFT, PUBLISHED, ARCHIVED
    calculation_method VARCHAR(50) NOT NULL, -- NET, GROSS
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS tbl_custom_form_version (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    version INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    UNIQUE (form_id, version)
);

CREATE UNIQUE INDEX uniq_active_form_version
ON tbl_custom_form_version (form_id)
WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS tbl_custom_form_field (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    form_version_id INTEGER NOT NULL REFERENCES tbl_custom_form_version(id),
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    start_date DATE NOT NULL,
    end_date DATE NULL,
    label VARCHAR(255) NOT NULL,
    section_type_id INTEGER NOT NULL REFERENCES tbl_section_type(id),
    description TEXT NULL,
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    coa_id VARCHAR(40) NOT NULL REFERENCES tbl_account(id),

    placeholder VARCHAR(255) NULL,
    min_value NUMERIC(14,2) NULL,
    max_value NUMERIC(14,2) NULL,
    field_order INTEGER NOT NULL,
    tax_type_id INTEGER NULL REFERENCES tbl_tax_type(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS tbl_custom_form_field_config (
    id VARCHAR(40) PRIMARY KEY,
    form_field_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form_field(id),
    tax_type_id INTEGER NOT NULL REFERENCES tbl_tax_type(id),
    arrangement_id VARCHAR(40) NULL REFERENCES tbl_arrangement(id),
    is_formula BOOLEAN NOT NULL DEFAULT FALSE,
    operator VARCHAR(10) NOT NULL, -- +, -, *, /
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS tbl_custom_form_field_formula_source (
    id VARCHAR(40) PRIMARY KEY,
    field_config_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_field_config(id),
    source_field_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_field(id),
    source_role VARCHAR(20) NOT NULL, -- PRIMARY, SECONDARY
    source_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    UNIQUE (field_config_id, source_field_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_field_formula_source;
DROP TABLE IF EXISTS tbl_custom_form_field_config;
DROP TABLE IF EXISTS tbl_custom_form_field;
DROP TABLE IF EXISTS tbl_custom_form_version;
DROP TABLE IF EXISTS tbl_custom_form;
DROP TABLE IF EXISTS tbl_tax_type;
DROP TABLE IF EXISTS tbl_section_type;

-- +goose StatementEnd
