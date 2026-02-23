-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_year (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),

    -- FY identity
    fy_label VARCHAR(9) NOT NULL,          -- e.g. '2024-25'
    start_date DATE NOT NULL,              -- e.g. 2024-07-01
    end_date DATE NOT NULL,                -- e.g. 2025-06-30

    -- Status
    is_current BOOLEAN DEFAULT FALSE,
    is_closed BOOLEAN DEFAULT FALSE,

    -- Controls
    closed_at TIMESTAMPTZ NULL,
    closed_by VARCHAR(40) NULL REFERENCES tbl_user(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_financial_year_clinic UNIQUE (clinic_id, start_date),
    CONSTRAINT chk_financial_year_dates CHECK (start_date < end_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_year;
-- +goose StatementEnd
