-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_auth_provider (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,

    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),

    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP NULL
);

-- ----------------------------------
-- UNIQUE CONSTRAINTS (SOFT DELETE SAFE)
-- ----------------------------------

-- One provider account → one user (globally)
CREATE UNIQUE INDEX ux_auth_provider_identity_active
ON tbl_auth_provider (provider, provider_user_id)
WHERE deleted_at IS NULL;

-- One provider per user (Google only once per user)
CREATE UNIQUE INDEX ux_auth_provider_user_provider_active
ON tbl_auth_provider (user_id, provider)
WHERE deleted_at IS NULL;

-- ----------------------------------
-- LOOKUP INDEXES
-- ----------------------------------

CREATE INDEX idx_auth_provider_user_active
ON tbl_auth_provider (user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_auth_provider_provider_active
ON tbl_auth_provider (provider)
WHERE deleted_at IS NULL;

CREATE INDEX idx_auth_provider_email_active
ON tbl_auth_provider (provider_email)
WHERE deleted_at IS NULL;

CREATE TRIGGER trg_auth_provider_updated_at
BEFORE UPDATE ON tbl_auth_provider
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_auth_provider;
-- +goose StatementEnd
