-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_transaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    source_form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    field_id VARCHAR(64) NULL,
    coa_id UUID NOT NULL REFERENCES tbl_account(id) ON DELETE RESTRICT,
    account_code VARCHAR(10) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    tax_category VARCHAR(32) NOT NULL,
    transaction_date DATE NOT NULL,
    reference VARCHAR(64) NOT NULL DEFAULT '',
    details TEXT NOT NULL DEFAULT '',
    gross_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'posted',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_transaction_clinic ON tbl_transaction(clinic_id);
CREATE INDEX IF NOT EXISTS idx_transaction_entry ON tbl_transaction(source_entry_id);
CREATE INDEX IF NOT EXISTS idx_transaction_date ON tbl_transaction(transaction_date);
CREATE INDEX IF NOT EXISTS idx_transaction_status ON tbl_transaction(status);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_transaction;
-- +goose StatementEnd
