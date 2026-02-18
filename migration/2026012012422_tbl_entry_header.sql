-- +goose Up
-- +goose StatementBegin

-- Add form_version_id to tbl_custom_form_entry
ALTER TABLE tbl_custom_form_entry ADD COLUMN IF NOT EXISTS form_version_id UUID REFERENCES tbl_custom_form_version(id) ON DELETE CASCADE;
UPDATE tbl_custom_form_entry cfe
SET form_version_id = f.form_version_id
FROM tbl_custom_form_field f
WHERE f.id = cfe.tbl_custom_form_field_id AND cfe.form_version_id IS NULL;
ALTER TABLE tbl_custom_form_entry ALTER COLUMN form_version_id SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tbl_custom_form_entry DROP COLUMN IF EXISTS form_version_id;

-- +goose StatementEnd
