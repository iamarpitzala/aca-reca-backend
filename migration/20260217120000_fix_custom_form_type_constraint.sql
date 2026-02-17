-- +goose Up
-- +goose StatementBegin

-- Fix form_type check constraint to allow 'BOTH' instead of 'REDUCTION'
-- Form types are: INCOME, EXPENSE, BOTH
-- REDUCTION is a section type for fields, not a form type

ALTER TABLE tbl_custom_form 
DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;

ALTER TABLE tbl_custom_form 
ADD CONSTRAINT tbl_custom_form_form_type_check 
CHECK (form_type IN ('INCOME', 'EXPENSE', 'BOTH'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tbl_custom_form 
DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;

ALTER TABLE tbl_custom_form 
ADD CONSTRAINT tbl_custom_form_form_type_check 
CHECK (form_type IN ('INCOME', 'EXPENSE', 'REDUCTION'));

-- +goose StatementEnd
