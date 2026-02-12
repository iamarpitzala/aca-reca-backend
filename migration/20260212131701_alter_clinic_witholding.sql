-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    ADD COLUMN with_holding_tax BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    DROP COLUMN IF EXISTS with_holding_tax;
-- +goose StatementEnd
