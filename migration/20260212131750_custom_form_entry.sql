-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_custom_form_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    value NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_by UUID NOT NULL REFERENCES tbl_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_custom_form_entry;

-- +goose StatementEnd
