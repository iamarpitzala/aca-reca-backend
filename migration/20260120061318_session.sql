-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,

    refresh_token VARCHAR(255) NOT NULL,

    user_agent TEXT,
    ip_address VARCHAR(45),

    expires_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP NULL
);

-- -------------------------
-- UNIQUE (ACTIVE SESSIONS ONLY)
-- -------------------------
CREATE UNIQUE INDEX ux_session_refresh_token_active
ON tbl_session (refresh_token)
WHERE deleted_at IS NULL;

-- -------------------------
-- COMMON LOOKUPS
-- -------------------------
CREATE INDEX idx_session_user_active
ON tbl_session (user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_session_expires_active
ON tbl_session (expires_at)
WHERE deleted_at IS NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_session;
-- +goose StatementEnd
