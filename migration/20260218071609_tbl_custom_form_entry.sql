-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_custom_form_entry (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    form_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form(id),
    form_version_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form_version(id),
    field_id VARCHAR(40) NOT NULL REFERENCES tbl_custom_form_field(id),
    value NUMERIC(14,2) NOT NULL,
    gst_amount NUMERIC(14,2) NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_entry;
-- +goose StatementEnd
