-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_custom_form
    ADD COLUMN IF NOT EXISTS outwork_enabled BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS outwork_rate_percent NUMERIC(5,2) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_custom_form
    DROP COLUMN IF EXISTS outwork_enabled,
    DROP COLUMN IF EXISTS outwork_rate_percent;
-- +goose StatementEnd
