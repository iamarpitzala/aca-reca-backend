-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    ADD COLUMN share_type VARCHAR(50) NOT NULL DEFAULT 'PERCENTAGE',
    ADD COLUMN clinic_share INT NOT NULL DEFAULT 50,
    ADD COLUMN owner_share INT NOT NULL DEFAULT 50;

ALTER TABLE tbl_clinic
    ADD CONSTRAINT chk_share_type
    CHECK (share_type IN ('FIXED', 'PERCENTAGE'));

ALTER TABLE tbl_clinic
    ADD CONSTRAINT chk_share_percentage
    CHECK (
        share_type != 'PERCENTAGE'
        OR (clinic_share + owner_share = 100)
    );

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

