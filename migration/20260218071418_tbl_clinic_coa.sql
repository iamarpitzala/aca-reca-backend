-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_clinic_coa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    coa_id UUID NOT NULL REFERENCES tbl_account(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- One account per clinic (active rows only)
CREATE UNIQUE INDEX IF NOT EXISTS ux_clinic_coa_clinic_account_active
ON tbl_clinic_coa (clinic_id, coa_id)
WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_clinic_coa_clinic_account_active;
DROP TABLE IF EXISTS tbl_clinic_coa;
-- +goose StatementEnd
