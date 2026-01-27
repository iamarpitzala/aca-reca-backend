-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_quarter(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

ALTER TABLE tbl_quarter
ADD CONSTRAINT chk_quarter_date_range
CHECK (start_date < end_date);

INSERT INTO tbl_quarter (name, start_date, end_date)
VALUES
    ('Q1 2025', DATE '2025-01-01', DATE '2025-03-31'),
    ('Q2 2025', DATE '2025-04-01', DATE '2025-06-30'),
    ('Q3 2025', DATE '2025-07-01', DATE '2025-09-30'),
    ('Q4 2025', DATE '2025-10-01', DATE '2025-12-31');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_quarter;
-- +goose StatementEnd
