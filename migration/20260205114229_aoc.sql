-- +goose Up
-- +goose StatementBegin

-- =========================
-- Account Type Master
-- =========================
CREATE TABLE IF NOT EXISTS tbl_account_type (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- =========================
-- Account Tax Master
-- =========================
CREATE TABLE IF NOT EXISTS tbl_account_tax (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    rate NUMERIC(5,2) NOT NULL DEFAULT 0, -- e.g. 10.00
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- =========================
-- Seed Account Types (Idempotent)
-- =========================
INSERT INTO tbl_account_type (name, description)
VALUES
    ('Asset', 'Asset account type'),
    ('Liability', 'Liability account type'),
    ('Equity', 'Equity account type'),
    ('Revenue', 'Revenue account type'),
    ('Expense', 'Expense account type')
ON CONFLICT (name) DO NOTHING;

-- =========================
-- Seed Tax Types (Xero-style)
-- =========================
INSERT INTO tbl_account_tax (name, rate, description)
VALUES
    ('Tax Exempt', 0.00, 'No tax applied'),
    ('Tax on Purchases', 0.00, 'GST/VAT on purchases'),
    ('Tax on Sales', 0.00, 'GST/VAT on sales')
ON CONFLICT (name) DO NOTHING;

-- =========================
-- Chart of Accounts
-- =========================
CREATE TABLE IF NOT EXISTS tbl_account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_type_id SMALLINT NOT NULL
        REFERENCES tbl_account_type(id),

    account_tax_id SMALLINT NOT NULL
        REFERENCES tbl_account_tax(id),

    code VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uq_account_code UNIQUE (code),
    CONSTRAINT uq_account_name UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_account_type
    ON tbl_account(account_type_id);

CREATE INDEX IF NOT EXISTS idx_account_tax
    ON tbl_account(account_tax_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_account;
DROP TABLE IF EXISTS tbl_account_tax;
DROP TABLE IF EXISTS tbl_account_type;
-- +goose StatementEnd
