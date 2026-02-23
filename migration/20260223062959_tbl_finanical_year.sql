-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_financial_year (
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

INSERT INTO tbl_financial_year (fy_label, start_date, end_date) VALUES
    ('2024-25', '2024-07-01', '2025-06-30'),
    ('2025-26', '2025-07-01', '2026-06-30'),
    ('2026-27', '2026-07-01', '2027-06-30'),
    ('2027-28', '2027-07-01', '2028-06-30'),
    ('2028-29', '2028-07-01', '2029-06-30'),
    ('2029-30', '2029-07-01', '2030-06-30')
ON CONFLICT (fy_label) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_year;
-- +goose StatementEnd
