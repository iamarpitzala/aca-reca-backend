-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_time_zone (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    timezone VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_timezone UNIQUE (timezone)
);

INSERT INTO tbl_time_zone (timezone) VALUES
    ('Australia/Sydney'),
    ('Australia/Melbourne'),
    ('Australia/Brisbane'),
    ('Australia/Perth'),
    ('Australia/Adelaide'),
    ('Australia/Hobart'),
    ('Australia/Darwin')
ON CONFLICT (timezone) DO NOTHING;

CREATE TABLE IF NOT EXISTS tbl_financial_setting (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),

    financial_year_id INTEGER NOT NULL REFERENCES tbl_financial_year(id),
    financial_quarter_id INTEGER NOT NULL REFERENCES tbl_financial_quarter(id),
    accounting_method VARCHAR(20) NOT NULL DEFAULT 'ACCRUAL' CHECK (accounting_method IN ('CASH', 'ACCRUAL')),
    
    tax_calculation_method VARCHAR(20) NOT NULL DEFAULT 'QUARTERLY' CHECK (tax_calculation_method IN ('QUARTERLY', 'ANNUALLY')),

    tax_sale_calculation_method VARCHAR(20) NOT NULL DEFAULT 'INCLUSIVE' CHECK (tax_sale_calculation_method IN ('INCLUSIVE', 'EXCLUSIVE')),

    tax_purchase_calculation_method VARCHAR(20) NOT NULL DEFAULT 'INCLUSIVE' CHECK (tax_purchase_calculation_method IN ('INCLUSIVE', 'EXCLUSIVE')),
    lock_date DATE,

    time_zone_id INTEGER NOT NULL REFERENCES tbl_time_zone(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_financial_setting;
DROP TABLE IF EXISTS tbl_time_zone;

-- +goose StatementEnd
