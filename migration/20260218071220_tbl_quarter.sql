-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_financial_year (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_financial_year_clinic
        FOREIGN KEY (clinic_id)
        REFERENCES tbl_clinic(id),
    CONSTRAINT chk_financial_year_date_range
        CHECK (start_date < end_date)
);

CREATE TABLE IF NOT EXISTS tbl_quarter (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    financial_year_id VARCHAR(40) NOT NULL REFERENCES tbl_financial_year(id),
    name VARCHAR(5) NOT NULL, -- Q1, Q2, Q3, Q4
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_quarter_clinic
        FOREIGN KEY (clinic_id)
        REFERENCES tbl_clinic(id),

    CONSTRAINT fk_quarter_financial_year
        FOREIGN KEY (financial_year_id)
        REFERENCES tbl_financial_year(id),

    CONSTRAINT chk_quarter_date_range
        CHECK (start_date < end_date)
);



-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_quarter;
DROP TABLE IF EXISTS tbl_financial_year;

-- +goose StatementEnd