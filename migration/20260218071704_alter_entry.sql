-- +goose Up
-- +goose StatementBegin

-- Add form_version_id to tbl_custom_form_entry
ALTER TABLE tbl_custom_form_entry ADD COLUMN IF NOT EXISTS form_version_id UUID REFERENCES tbl_custom_form_version(id) ON DELETE CASCADE;
UPDATE tbl_custom_form_entry cfe
SET form_version_id = f.form_version_id
FROM tbl_custom_form_field f
WHERE f.id = cfe.tbl_custom_form_field_id AND cfe.form_version_id IS NULL;
ALTER TABLE tbl_custom_form_entry ALTER COLUMN form_version_id SET NOT NULL;

ALTER TABLE tbl_entry_net_details
ADD COLUMN IF NOT EXISTS net_amount NUMERIC(14,2) NOT NULL DEFAULT 0;

ALTER TABLE tbl_custom_form_field
    DROP COLUMN IF EXISTS gst_type;

ALTER TABLE tbl_custom_form_field
    ADD COLUMN gst_type VARCHAR(50)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tbl_custom_form_entry DROP COLUMN IF EXISTS form_version_id;
ALTER TABLE tbl_entry_net_details DROP COLUMN IF EXISTS net_amount;
ALTER TABLE tbl_custom_form_field
    DROP COLUMN IF EXISTS gst_type;

ALTER TABLE tbl_custom_form_field
    ADD COLUMN gst_type VARCHAR(50)
-- +goose StatementEnd
