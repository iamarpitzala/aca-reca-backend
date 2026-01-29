-- +goose Up
-- +goose StatementBegin

ALTER TABLE tbl_financial_form
ADD COLUMN quarter_id UUID;

ALTER TABLE tbl_financial_form
ADD CONSTRAINT fk_tbl_financial_form_quarter
FOREIGN KEY (quarter_id)
REFERENCES tbl_quarter(id)
ON DELETE CASCADE;

CREATE INDEX idx_tbl_financial_form_quarter_id
ON tbl_financial_form(quarter_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tbl_financial_form_quarter_id;

ALTER TABLE tbl_financial_form
DROP CONSTRAINT IF EXISTS fk_tbl_financial_form_quarter;

ALTER TABLE tbl_financial_form
DROP COLUMN IF EXISTS quarter_id;
-- +goose StatementEnd
