-- +goose Up
-- +goose StatementBegin

ALTER TABLE tbl_session
ALTER COLUMN refresh_token TYPE TEXT;

DROP INDEX IF EXISTS idx_sessions_refresh_token;

CREATE UNIQUE INDEX IF NOT EXISTS ux_session_refresh_token_active
ON tbl_session (refresh_token)
WHERE deleted_at IS NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- ⚠️ Down migration is intentionally conservative
-- Reverting token length may FAIL if long tokens exist

DROP INDEX IF EXISTS ux_session_refresh_token_active;

ALTER TABLE tbl_session
ALTER COLUMN refresh_token TYPE VARCHAR(255);

CREATE UNIQUE INDEX idx_sessions_refresh_token
ON tbl_session (refresh_token);

-- +goose StatementEnd
