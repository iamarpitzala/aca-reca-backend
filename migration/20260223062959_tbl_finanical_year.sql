-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_year (
    id SERIAL PRIMARY KEY,
    
    fy_label VARCHAR(9) NOT NULL UNIQUE,        -- e.g. 2024-25
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT chk_fy_dates CHECK (start_date < end_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_year;
-- +goose StatementEnd
