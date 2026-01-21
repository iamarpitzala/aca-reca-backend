-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON tbl_session(user_id);
CREATE INDEX idx_sessions_expires_at ON tbl_session(expires_at);
CREATE INDEX idx_sessions_deleted_at ON tbl_session(deleted_at);
CREATE INDEX idx_sessions_refresh_token ON tbl_session(refresh_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_session;
-- +goose StatementEnd
