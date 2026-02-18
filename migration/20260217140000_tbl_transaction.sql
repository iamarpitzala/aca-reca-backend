-- +goose Up
-- +goose StatementBegin

-- ============================================================
-- Transaction Header
-- One row per accounting transaction (e.g. an entry posting).
-- Scoped to a clinic for multi-tenant isolation.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_transaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    clinic_id UUID NOT NULL
        REFERENCES tbl_clinic(id) ON DELETE CASCADE,

    source_entry_id UUID
        REFERENCES tbl_custom_form_entry(id) ON DELETE SET NULL,

    reference_number VARCHAR(50),
    description TEXT,

    transaction_date DATE NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'POSTED', 'VOIDED')),

    created_by UUID NOT NULL
        REFERENCES tbl_user(id),

    posted_at TIMESTAMPTZ,
    voided_at TIMESTAMPTZ,
    void_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- ============================================================
-- Transaction Ledger (line items)
-- Each row is one side of the double-entry posting against a
-- Chart-of-Accounts (COA) account.
-- net_amount is signed: positive = debit, negative = credit.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_transaction_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    transaction_id UUID NOT NULL
        REFERENCES tbl_transaction(id) ON DELETE CASCADE,

    coa_id UUID NOT NULL
        REFERENCES tbl_account(id) ON DELETE RESTRICT,

    entry_type VARCHAR(10) NOT NULL
        CHECK (entry_type IN ('DEBIT', 'CREDIT')),

    amount NUMERIC(14,2) NOT NULL DEFAULT 0
        CHECK (amount >= 0),

    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,

    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,

    transaction_date DATE NOT NULL,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- Indexes
-- ============================================================

-- tbl_transaction
CREATE INDEX IF NOT EXISTS idx_transaction_clinic_id
    ON tbl_transaction(clinic_id);

CREATE INDEX IF NOT EXISTS idx_transaction_clinic_date
    ON tbl_transaction(clinic_id, transaction_date);

CREATE INDEX IF NOT EXISTS idx_transaction_status
    ON tbl_transaction(status);

CREATE INDEX IF NOT EXISTS idx_transaction_source_entry
    ON tbl_transaction(source_entry_id)
    WHERE source_entry_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_transaction_created_by
    ON tbl_transaction(created_by);

CREATE INDEX IF NOT EXISTS idx_transaction_active
    ON tbl_transaction(id)
    WHERE deleted_at IS NULL;

-- tbl_transaction_ledger
CREATE INDEX IF NOT EXISTS idx_txn_ledger_transaction_id
    ON tbl_transaction_ledger(transaction_id);

CREATE INDEX IF NOT EXISTS idx_txn_ledger_coa_id
    ON tbl_transaction_ledger(coa_id);

CREATE INDEX IF NOT EXISTS idx_txn_ledger_date
    ON tbl_transaction_ledger(transaction_date);

CREATE INDEX IF NOT EXISTS idx_txn_ledger_coa_date
    ON tbl_transaction_ledger(coa_id, transaction_date);

-- ============================================================
-- Triggers
-- ============================================================
CREATE TRIGGER trg_transaction_updated_at
    BEFORE UPDATE ON tbl_transaction
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_transaction_updated_at ON tbl_transaction;

DROP INDEX IF EXISTS idx_txn_ledger_coa_date;
DROP INDEX IF EXISTS idx_txn_ledger_date;
DROP INDEX IF EXISTS idx_txn_ledger_coa_id;
DROP INDEX IF EXISTS idx_txn_ledger_transaction_id;

DROP INDEX IF EXISTS idx_transaction_active;
DROP INDEX IF EXISTS idx_transaction_created_by;
DROP INDEX IF EXISTS idx_transaction_source_entry;
DROP INDEX IF EXISTS idx_transaction_status;
DROP INDEX IF EXISTS idx_transaction_clinic_date;
DROP INDEX IF EXISTS idx_transaction_clinic_id;

DROP TABLE IF EXISTS tbl_transaction_ledger;
DROP TABLE IF EXISTS tbl_transaction;

-- +goose StatementEnd
