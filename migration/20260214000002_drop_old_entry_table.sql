-- +goose Up
-- +goose StatementBegin

-- Drop old tbl_custom_form_entry table and its indexes
-- Note: This should only be run after data migration is complete
-- and all references have been updated

-- First, drop the foreign key constraint from tbl_entry_header.original_entry_id
ALTER TABLE tbl_entry_header 
    DROP CONSTRAINT IF EXISTS tbl_entry_header_original_entry_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS idx_custom_form_entry_deleted_at;
DROP INDEX IF EXISTS idx_custom_form_entry_entry_date;
DROP INDEX IF EXISTS idx_custom_form_entry_quarter_id;
DROP INDEX IF EXISTS idx_custom_form_entry_clinic_id;
DROP INDEX IF EXISTS idx_custom_form_entry_form_id;

-- Now drop the table
DROP TABLE IF EXISTS tbl_custom_form_entry;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Recreate the old table structure (for rollback)
CREATE TABLE IF NOT EXISTS tbl_custom_form_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    form_name VARCHAR(255) NOT NULL,
    form_type VARCHAR(50) NOT NULL CHECK (form_type IN ('income', 'expense', 'both')),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    quarter_id UUID REFERENCES tbl_quarter(id) ON DELETE SET NULL,
    values JSONB NOT NULL DEFAULT '[]',
    calculations JSONB NOT NULL DEFAULT '{}',
    entry_date DATE NOT NULL,
    description TEXT,
    remarks TEXT,
    payment_responsibility VARCHAR(50) CHECK (payment_responsibility IN ('owner', 'clinic')),
    deductions JSONB,
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_custom_form_entry_form_id ON tbl_custom_form_entry(form_id);
CREATE INDEX idx_custom_form_entry_clinic_id ON tbl_custom_form_entry(clinic_id);
CREATE INDEX idx_custom_form_entry_quarter_id ON tbl_custom_form_entry(quarter_id);
CREATE INDEX idx_custom_form_entry_entry_date ON tbl_custom_form_entry(entry_date);
CREATE INDEX idx_custom_form_entry_deleted_at ON tbl_custom_form_entry(deleted_at);

-- Recreate the foreign key constraint from tbl_entry_header.original_entry_id
ALTER TABLE tbl_entry_header
    ADD CONSTRAINT tbl_entry_header_original_entry_id_fkey
    FOREIGN KEY (original_entry_id)
    REFERENCES tbl_custom_form_entry(id)
    ON DELETE SET NULL;

-- +goose StatementEnd
