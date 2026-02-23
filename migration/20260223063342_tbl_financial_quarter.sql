-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_financial_quarter (
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
    CONSTRAINT chk_quarter_dates CHECK (start_date < end_date),
);

INSERT INTO tbl_financial_quarter (financial_year_id, name, quarter_number, start_date, end_date) VALUES
    (1, 'Q1', 1, '2024-07-01', '2024-09-30'),
    (1, 'Q2', 2, '2024-10-01', '2024-12-31'),
    (1, 'Q3', 3, '2025-01-01', '2025-03-31'),
    (1, 'Q4', 4, '2025-04-01', '2025-06-30'),
    (2, 'Q1', 1, '2025-07-01', '2025-09-30'),
    (2, 'Q2', 2, '2025-10-01', '2025-12-31'),
    (2, 'Q3', 3, '2026-01-01', '2026-03-31'),
    (2, 'Q4', 4, '2026-04-01', '2026-06-30'),
    (3, 'Q1', 1, '2026-07-01', '2026-09-30'),
    (3, 'Q2', 2, '2026-10-01', '2026-12-31'),
    (3, 'Q3', 3, '2027-01-01', '2027-03-31'),
    (3, 'Q4', 4, '2027-04-01', '2027-06-30'),
    (4, 'Q1', 1, '2028-07-01', '2028-09-30'),
    (4, 'Q2', 2, '2028-10-01', '2028-12-31'),
    (4, 'Q3', 3, '2029-01-01', '2029-03-31'),
    (4, 'Q4', 4, '2029-04-01', '2029-06-30'),
    (5, 'Q1', 1, '2029-07-01', '2029-09-30'),
    (5, 'Q2', 2, '2029-10-01', '2029-12-31'),
    (5, 'Q3', 3, '2030-01-01', '2030-03-31'),
    (5, 'Q4', 4, '2030-04-01', '2030-06-30'),
    (6, 'Q1', 1, '2030-07-01', '2030-09-30'),
    (6, 'Q2', 2, '2030-10-01', '2030-12-31'),
    (6, 'Q3', 3, '2031-01-01', '2031-03-31'),
    (6, 'Q4', 4, '2031-04-01', '2031-06-30'),
ON CONFLICT (financial_year_id, quarter_number) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_quarter;
-- +goose StatementEnd
