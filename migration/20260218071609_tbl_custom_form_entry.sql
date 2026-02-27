-- +goose Up
-- +goose StatementBegin

CREATE TABLE tbl_custom_form_entry (
    id VARCHAR(40) PRIMARY KEY,
    form_version_id INTEGER NOT NULL
        REFERENCES tbl_custom_form_version(id),

    submitted_by VARCHAR(40) REFERENCES tbl_user(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);


CREATE TABLE tbl_custom_form_entry_value (
    id SERIAL PRIMARY KEY,

    entry_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,

    field_id VARCHAR(40) NOT NULL
        REFERENCES tbl_custom_form_field(id),

    value NUMERIC(14,2) NOT NULL,
    gst_amount NUMERIC(14,2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    UNIQUE (entry_id, field_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_entry_value;
DROP TABLE IF EXISTS tbl_custom_form_entry;
-- +goose StatementEnd
