-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_transaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    source_entry_id UUID NOT NULL REFERENCES tbl_custom_form_entry(id) ON DELETE CASCADE,
    source_form_id UUID NOT NULL REFERENCES tbl_custom_form(id) ON DELETE CASCADE,
    tbl_custom_form_field_id UUID NOT NULL REFERENCES tbl_custom_form_field(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'POSTED' CHECK (status IN ('POSTED', 'DRAFT', 'VOIDED')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_transaction;
-- +goose StatementEnd
