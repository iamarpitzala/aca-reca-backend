-- +goose Up
-- +goose StatementBegin

-- Restrict role to known values
ALTER TABLE tbl_user_clinic
DROP CONSTRAINT IF EXISTS chk_user_clinic_role;

ALTER TABLE tbl_user_clinic
ADD CONSTRAINT chk_user_clinic_role
CHECK (role IN ('OWNER', 'MEMBER', 'VIEWER'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tbl_user_clinic DROP CONSTRAINT IF EXISTS chk_user_clinic_role;

-- +goose StatementEnd
