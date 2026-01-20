-- +goose Up
-- +goose StatementBegin
-- Change refresh_token from VARCHAR(255) to TEXT to support longer JWT tokens
ALTER TABLE tbl_session 
ALTER COLUMN refresh_token TYPE TEXT;

-- Recreate the unique index on refresh_token (it might have been dropped)
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_refresh_token ON tbl_session(refresh_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert back to VARCHAR(255) - Note: This may fail if existing tokens are longer
DROP INDEX IF EXISTS idx_sessions_refresh_token;
ALTER TABLE tbl_session 
ALTER COLUMN refresh_token TYPE VARCHAR(255);
CREATE UNIQUE INDEX idx_sessions_refresh_token ON tbl_session(refresh_token);
-- +goose StatementEnd
