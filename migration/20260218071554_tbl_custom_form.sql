-- +goose Up
-- +goose StatementBegin

CREATE TABLE tbl_tax_type (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL, -- INCLUSIVE, EXCLUSIVE, MANUAL
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

INSERT INTO tbl_tax_type (code, name, description) VALUES
('INCLUSIVE', 'Inclusive', 'Inclusive tax'),
('EXCLUSIVE', 'Exclusive', 'Exclusive tax'),
('MANUAL', 'Manual', 'Manual tax');


CREATE TABLE tbl_custom_form (
    id VARCHAR(40) PRIMARY KEY,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL, -- DRAFT, PUBLISHED, ARCHIVED
    calculation_method VARCHAR(50) NOT NULL, -- NET, GROSS
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE tbl_custom_form_version (
    id SERIAL PRIMARY KEY,
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    version INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (form_id, version)
);

CREATE UNIQUE INDEX uniq_active_form_version
ON tbl_custom_form_version (form_id)
WHERE is_active = TRUE;

CREATE TABLE tbl_section_type (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL, -- COLLECTION, COST, SERVICE_FACILITY, OTHER_COST
    name VARCHAR(255) NOT NULL,
    description TEXT
);

INSERT INTO tbl_section_type (code, name, description) VALUES
    ('COLLECTION', 'Collection', 'Collection section'),
    ('COST', 'Cost', 'Cost section'),
    ('SERVICE_FACILITY', 'Service Facility', 'Service Facility section'),
    ('OTHER_COST', 'Other Cost', 'Other Cost section');


CREATE TABLE tbl_custom_form_section (
    id SERIAL PRIMARY KEY,
    form_version_id INTEGER NOT NULL REFERENCES tbl_custom_form_version(id),
    section_type_id INTEGER NOT NULL REFERENCES tbl_section_type(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    section_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE tbl_payment_responsibility (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL, -- OWNER, CLINIC
    name VARCHAR(255) NOT NULL
);
INSERT INTO tbl_payment_responsibility (code, name) VALUES
    ('OWNER', 'Owner'),
    ('CLINIC', 'Clinic');

CREATE TABLE tbl_custom_form_field (
    id VARCHAR(40) PRIMARY KEY,
    section_id INTEGER NOT NULL REFERENCES tbl_custom_form_section(id),
    label VARCHAR(255) NOT NULL,
    payment_responsibility_id SMALLINT NULL 
        REFERENCES tbl_payment_responsibility(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_field;
DROP TABLE IF EXISTS tbl_custom_form_version;
DROP TABLE IF EXISTS tbl_custom_form;
DROP TABLE IF EXISTS tbl_tax_type;
DROP TABLE IF EXISTS tbl_section_type;
DROP TABLE IF EXISTS tbl_payment_responsibility;

-- +goose StatementEnd
