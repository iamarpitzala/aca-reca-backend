-- +goose Up
-- +goose StatementBegin

-- Update transaction table to reference new tbl_entry_header instead of tbl_custom_form_entry
-- First, drop the foreign key constraint
ALTER TABLE tbl_transaction 
    DROP CONSTRAINT IF EXISTS tbl_transaction_source_entry_id_fkey;

-- Add new foreign key constraint referencing tbl_entry_header
ALTER TABLE tbl_transaction
    ADD CONSTRAINT tbl_transaction_source_entry_id_fkey 
    FOREIGN KEY (source_entry_id) 
    REFERENCES tbl_entry_header(id) 
    ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert back to old table reference
ALTER TABLE tbl_transaction 
    DROP CONSTRAINT IF EXISTS tbl_transaction_source_entry_id_fkey;

ALTER TABLE tbl_transaction
    ADD CONSTRAINT tbl_transaction_source_entry_id_fkey 
    FOREIGN KEY (source_entry_id) 
    REFERENCES tbl_custom_form_entry(id) 
    ON DELETE CASCADE;

-- +goose StatementEnd
