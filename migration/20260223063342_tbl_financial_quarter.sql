-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_quarter (
    id SERIAL PRIMARY KEY,
    financial_year_id INTEGER NOT NULL REFERENCES tbl_financial_year(id) ON DELETE CASCADE,

    quarter_number SMALLINT NOT NULL CHECK (quarter_number BETWEEN 1 AND 4),
    name VARCHAR(2) NOT NULL, -- Q1..Q4

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT uq_fy_quarter UNIQUE (financial_year_id, quarter_number),
    CONSTRAINT chk_quarter_dates CHECK (start_date < end_date)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_quarter;
-- +goose StatementEnd
