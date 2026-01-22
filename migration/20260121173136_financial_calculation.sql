-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_financial_calculation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    financial_form_id UUID NOT NULL REFERENCES tbl_financial_form(id) ON DELETE CASCADE,
    input_data JSONB NOT NULL,
    calculated_data JSONB NOT NULL,
    bas_mapping JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL
);

CREATE INDEX idx_financial_calculation_form_id ON tbl_financial_calculation(financial_form_id);
CREATE INDEX idx_financial_calculation_created_at ON tbl_financial_calculation(created_at);
CREATE INDEX idx_financial_calculation_created_by ON tbl_financial_calculation(created_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_financial_calculation;
-- +goose StatementEnd
