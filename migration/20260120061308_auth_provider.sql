-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_auth_provider (
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
    deleted_at TIMESTAMP
);

CREATE INDEX idx_oauth_providers_user_id ON tbl_auth_provider(user_id);
CREATE INDEX idx_oauth_providers_provider ON tbl_auth_provider(provider);
CREATE INDEX idx_oauth_providers_provider_user_id ON tbl_auth_provider(provider_user_id);
CREATE INDEX idx_oauth_providers_deleted_at ON tbl_auth_provider(deleted_at);
CREATE UNIQUE INDEX idx_oauth_providers_user_provider ON tbl_auth_provider(user_id, provider) WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_auth_provider;
-- +goose StatementEnd
