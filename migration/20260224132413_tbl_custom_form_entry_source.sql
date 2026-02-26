-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_custom_form_entry_source (
    id VARCHAR(40) PRIMARY KEY NOT NULL,

    entry_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_entry(id),

    field_config_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_field_config(id),

    source_field_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_field(id),

    source_value NUMERIC(14,2) NOT NULL,
    source_gst_amount NUMERIC(14,2) NULL,

    source_role VARCHAR(20) NOT NULL, -- PRIMARY, SECONDARY
    source_order INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (entry_id, source_field_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_entry_source;
-- +goose StatementEnd
