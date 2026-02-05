-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_account_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS tbl_account_tax (
    id SERIAL PRIMARY KEY,
    tbl_account_type_id INT NOT NULL REFERENCES tbl_account_type(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS tbl_aoc (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tbl_account_type_id INT NOT NULL REFERENCES tbl_account_type(id) ON DELETE CASCADE,
    tbl_account_tax_id INT NOT NULL REFERENCES tbl_account_tax(id) ON DELETE CASCADE,
    code VARCHAR(4) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

INSERT INTO tbl_account_type (name, description) VALUES
('Asset', 'Asset account type'),
('Liability', 'Liability account type'),
('Equity', 'Equity account type'),
('Revenue', 'Revenue account type'),
('Expense', 'Expense account type');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_aoc;
DROP TABLE IF EXISTS tbl_account_tax;
DROP TABLE IF EXISTS tbl_account_type;
-- +goose StatementEnd
