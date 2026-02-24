-- +goose Up
-- +goose StatementBegin
-- Prevent duplicate ABN among active clinics (same ABN allowed after soft delete).
CREATE UNIQUE INDEX ux_clinic_abn_active
ON tbl_clinic (abn_number)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_clinic_abn_active;
-- +goose StatementEnd
