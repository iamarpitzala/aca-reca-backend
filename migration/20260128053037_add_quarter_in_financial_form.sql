-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_financial_form
ADD CONSTRAINT fk_tbl_financial_form_quarter
FOREIGN KEY (quarter_id)
REFERENCES tbl_quarter(id)
ON DELETE CASCADE;



ALTER TABLE tbl_financial_form
VALIDATE CONSTRAINT fk_tbl_financial_form_quarter;

CREATE INDEX idx_tbl_financial_form_quarter_id
ON tbl_financial_form(quarter_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
