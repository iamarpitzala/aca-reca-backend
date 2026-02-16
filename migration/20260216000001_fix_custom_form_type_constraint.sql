-- +goose Up
-- +goose StatementBegin

-- Drop the existing check constraint
ALTER TABLE tbl_custom_form DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;

-- Add the correct check constraint that allows INCOME, EXPENSE, or BOTH
ALTER TABLE tbl_custom_form 
ADD CONSTRAINT tbl_custom_form_form_type_check 
CHECK (form_type IN ('INCOME', 'EXPENSE', 'BOTH'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert to the old constraint (for rollback purposes)
ALTER TABLE tbl_custom_form DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;

ALTER TABLE tbl_custom_form 
ADD CONSTRAINT tbl_custom_form_form_type_check 
CHECK (form_type IN ('INCOME', 'EXPENSE', 'REDUCTION'));

-- +goose StatementEnd
