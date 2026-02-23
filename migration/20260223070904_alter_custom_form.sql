-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_custom_form ADD COLUMN arrangement_id VARCHAR(40) NULL REFERENCES tbl_arrangement(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_custom_form DROP COLUMN IF EXISTS arrangement_id;
-- +goose StatementEnd
