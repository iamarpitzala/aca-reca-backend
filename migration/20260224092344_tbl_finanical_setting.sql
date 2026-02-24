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

-- Application uses tbl_clinic_financial_settings (see clinic_financial_setting adapter/domain).
CREATE TABLE IF NOT EXISTS tbl_clinic_financial_settings (
    id VARCHAR(40) PRIMARY KEY,

    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    clinic_financial_year_id INTEGER NOT NULL
        REFERENCES tbl_clinic_financial_year(id),

    accounting_method VARCHAR(20) NOT NULL
        CHECK (accounting_method IN ('CASH', 'ACCRUAL')),

    gst_registered BOOLEAN DEFAULT true,

    gst_reporting_frequency VARCHAR(20)
        CHECK (gst_reporting_frequency IN ('QUARTERLY', 'ANNUALLY')),

    default_amount_mode VARCHAR(20)
        CHECK (default_amount_mode IN ('INCLUSIVE', 'EXCLUSIVE')),

    lock_date DATE NULL,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_clinic_fy_settings UNIQUE (clinic_financial_year_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_clinic_financial_settings;
DROP TABLE IF EXISTS tbl_time_zone;

-- +goose StatementEnd
