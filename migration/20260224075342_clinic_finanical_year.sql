-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_clinic_financial_year (
    id SERIAL PRIMARY KEY,

    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    financial_year_id INTEGER NOT NULL REFERENCES tbl_financial_year(id),

    is_current BOOLEAN DEFAULT FALSE,
    is_closed BOOLEAN DEFAULT FALSE,

    closed_at TIMESTAMPTZ NULL,
    closed_by VARCHAR(40) NULL REFERENCES tbl_user(id),

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_clinic_fy UNIQUE (clinic_id, financial_year_id)
);

CREATE UNIQUE INDEX uq_clinic_current_fy
ON tbl_clinic_financial_year (clinic_id)
WHERE is_current = true AND deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_clinic_finanical_year;
-- +goose StatementEnd
