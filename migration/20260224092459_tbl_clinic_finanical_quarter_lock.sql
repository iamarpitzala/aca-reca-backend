-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_clinic_financial_quarter_lock (
    id SERIAL PRIMARY KEY,

    clinic_financial_year_id INTEGER NOT NULL
        REFERENCES tbl_clinic_financial_year(id),

    financial_quarter_id INTEGER NOT NULL
        REFERENCES tbl_financial_quarter(id),

    is_closed BOOLEAN DEFAULT FALSE,
    closed_at TIMESTAMPTZ NULL,
    closed_by VARCHAR(40) REFERENCES tbl_user(id),

    CONSTRAINT uq_clinic_quarter UNIQUE (
        clinic_financial_year_id,
        financial_quarter_id
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_clinic_financial_quarter_lock;
-- +goose StatementEnd
