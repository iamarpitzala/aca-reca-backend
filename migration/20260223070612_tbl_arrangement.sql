-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_arrangement (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    method VARCHAR(50) NOT NULL DEFAULT 'GROSS' CHECK (method IN ('GROSS', 'NET')),
    name VARCHAR(255) NOT NULL,
    percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    amount NUMERIC(14,2) NULL,
    description TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_arrangement;
-- +goose StatementEnd
