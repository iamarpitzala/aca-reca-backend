-- +goose Up
-- +goose StatementBegin
   -- Find and drop the old constraint
   ALTER TABLE tbl_custom_form DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;
   
   -- Add the new constraint
   ALTER TABLE tbl_custom_form ADD CONSTRAINT tbl_custom_form_form_type_check 
       CHECK (form_type IN ('INCOME', 'EXPENSE', 'REDUCTION', 'BOTH'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
   ALTER TABLE tbl_custom_form DROP CONSTRAINT IF EXISTS tbl_custom_form_form_type_check;
   ALTER TABLE tbl_custom_form ADD CONSTRAINT tbl_custom_form_form_type_check 
       CHECK (form_type IN ('INCOME', 'EXPENSE', 'BOTH'));
-- +goose StatementEnd
