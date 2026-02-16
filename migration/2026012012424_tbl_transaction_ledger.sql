-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_transaction_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES tbl_transaction(id) ON DELETE CASCADE,
    coa_id UUID NOT NULL REFERENCES tbl_account(id) ON DELETE RESTRICT,
    account_code VARCHAR(10) NOT NULL DEFAULT '',
    account_name VARCHAR(255) NOT NULL DEFAULT '',
    tax_category VARCHAR(50) NOT NULL DEFAULT '',
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reference TEXT NOT NULL DEFAULT '',
    details TEXT NOT NULL DEFAULT '',
    gross_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_transaction_id ON tbl_transaction_ledger(transaction_id);
CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_coa_id ON tbl_transaction_ledger(coa_id);
CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_account_code ON tbl_transaction_ledger(account_code);
CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_account_name ON tbl_transaction_ledger(account_name);
CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_tax_category ON tbl_transaction_ledger(tax_category);
CREATE INDEX IF NOT EXISTS idx_tbl_transaction_ledger_transaction_date ON tbl_transaction_ledger(transaction_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_tbl_transaction_ledger_transaction_id;
DROP INDEX IF EXISTS idx_tbl_transaction_ledger_coa_id;
DROP INDEX IF EXISTS idx_tbl_transaction_ledger_account_code;
DROP INDEX IF EXISTS idx_tbl_transaction_ledger_account_name;
DROP INDEX IF EXISTS idx_tbl_transaction_ledger_tax_category;
DROP INDEX IF EXISTS idx_tbl_transaction_ledger_transaction_date;

DROP TABLE IF EXISTS tbl_transaction_ledger;

-- +goose StatementEnd
