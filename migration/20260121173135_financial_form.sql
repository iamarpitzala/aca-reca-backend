-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    calculation_method VARCHAR(50) NOT NULL CHECK (calculation_method IN ('net', 'gross')),
    configuration JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_financial_form_clinic_id ON tbl_financial_form(clinic_id);
CREATE INDEX idx_financial_form_method ON tbl_financial_form(calculation_method);
CREATE INDEX idx_financial_form_deleted_at ON tbl_financial_form(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_form;
-- +goose StatementEnd
