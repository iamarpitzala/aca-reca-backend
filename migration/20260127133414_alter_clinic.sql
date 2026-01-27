-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    ADD COLUMN share_type VARCHAR(50) NOT NULL DEFAULT 'FIXED',
    ADD COLUMN clinic_share INT NOT NULL DEFAULT 0,
    ADD COLUMN owner_share INT NOT NULL DEFAULT 0;

-- Optional but highly recommended guardrails
ALTER TABLE tbl_clinic
    ADD CONSTRAINT chk_share_percentage
    CHECK (clinic_share + owner_share = 100);

ALTER TABLE tbl_clinic
    ADD CONSTRAINT chk_share_type
    CHECK (share_type IN ('FIXED', 'PERCENTAGE'));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    DROP CONSTRAINT IF EXISTS chk_share_percentage,
    DROP CONSTRAINT IF EXISTS chk_share_type;

ALTER TABLE tbl_clinic
    DROP COLUMN IF EXISTS share_type,
    DROP COLUMN IF EXISTS clinic_share,
    DROP COLUMN IF EXISTS owner_share;
-- +goose StatementEnd

