-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS tbl_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(50),
    avatar_url TEXT,

    is_email_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP NULL
);

-- -------------------------
-- UNIQUE (ACTIVE USERS ONLY)
-- -------------------------
CREATE UNIQUE INDEX ux_user_email_active
ON tbl_user (LOWER(email))
WHERE deleted_at IS NULL;

-- -------------------------
-- COMMON LOOKUPS
-- -------------------------
CREATE INDEX idx_user_active
ON tbl_user (id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_user_email_lookup
ON tbl_user (LOWER(email))
WHERE deleted_at IS NULL;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_user_updated_at
BEFORE UPDATE ON tbl_user
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_user;
-- +goose StatementEnd
