-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_custom_form ADD COLUMN arrangement_id VARCHAR(40) NULL REFERENCES tbl_arrangement(id);

CREATE UNIQUE INDEX ux_clinic_abn_active
ON tbl_clinic (abn_number)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_custom_form DROP COLUMN IF EXISTS arrangement_id;
DROP INDEX IF EXISTS ux_clinic_abn_active;
-- +goose StatementEnd
