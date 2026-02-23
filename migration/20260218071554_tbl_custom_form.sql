-- +goose Up
-- +goose StatementBegin

CREATE TABLE tbl_tax_type (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCLUSIVE', 'EXCLUSIVE', 'MANUAL')),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

INSERT INTO tbl_tax_type (name, type, description) VALUES
    ('Inclusive', 'INCLUSIVE', 'Inclusive tax type'),
    ('Exclusive', 'EXCLUSIVE', 'Exclusive tax type'),
    ('Manual', 'MANUAL', 'Manual tax type')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE tbl_section_type (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

INSERT INTO tbl_section_type (name, type, description) VALUES
    ('INCOME', 'INCOME', 'Income section'),
    ('EXPENSE', 'EXPENSE', 'Expense section')
ON CONFLICT (name) DO NOTHING;


CREATE TABLE tbl_custom_form (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED')),
    calculation_method VARCHAR(50) NOT NULL CHECK (calculation_method IN ('NET', 'GROSS')) DEFAULT 'NET',
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE tbl_custom_form_version (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    version INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    UNIQUE (form_id, version)
);

CREATE TABLE tbl_custom_form_field (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    form_version_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form_version(id),
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    label VARCHAR(255) NOT NULL,
    section_type_id VARCHAR(40) NOT NULL REFERENCES tbl_section_type(id),
    description TEXT NULL,
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    coa_id VARCHAR(40) NOT NULL REFERENCES tbl_account(id),

    placeholder VARCHAR(255) NULL,
    min_value NUMERIC(14,2) NULL,
    max_value NUMERIC(14,2) NULL,
    field_order INTEGER NOT NULL,
    tax_type_id VARCHAR(40) NULL REFERENCES tbl_tax_type(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_field;
DROP TABLE IF EXISTS tbl_custom_form_version;
DROP TABLE IF EXISTS tbl_custom_form;
DROP TABLE IF EXISTS tbl_tax_type;
DROP TABLE IF EXISTS tbl_section_type;
DROP TABLE IF EXISTS tbl_tax_method;
-- +goose StatementEnd
