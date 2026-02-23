-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_quarter (
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    financial_year_id INTEGER NOT NULL REFERENCES tbl_financial_year(id),
    name VARCHAR(5) NOT NULL, -- Q1, Q2, Q3, Q4

    quarter_number SMALLINT NOT NULL CHECK (quarter_number BETWEEN 1 AND 4),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    is_closed BOOLEAN DEFAULT FALSE,
    closed_at TIMESTAMPTZ NULL,
    closed_by VARCHAR(40) NULL REFERENCES tbl_user(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_fy_quarter UNIQUE (financial_year_id, quarter_number),
    CONSTRAINT chk_quarter_dates CHECK (start_date < end_date)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_quarter;
-- +goose StatementEnd
