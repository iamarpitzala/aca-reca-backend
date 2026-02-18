-- +goose Up
-- +goose StatementBegin

ALTER TABLE tbl_clinic
    ADD COLUMN IF NOT EXISTS share_type VARCHAR(50) NOT NULL DEFAULT 'PERCENTAGE',
    ADD COLUMN IF NOT EXISTS clinic_share INT NOT NULL DEFAULT 50,
    ADD COLUMN IF NOT EXISTS owner_share INT NOT NULL DEFAULT 50,
    ADD COLUMN IF NOT EXISTS method_type VARCHAR(50) NOT NULL DEFAULT 'NET'
        CHECK (method_type IN ('NET', 'GROSS')),
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS with_holding_tax BOOLEAN NOT NULL DEFAULT FALSE;

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
    DROP COLUMN IF EXISTS with_holding_tax,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS method_type,
    DROP COLUMN IF EXISTS share_type,
    DROP COLUMN IF EXISTS clinic_share,
    DROP COLUMN IF EXISTS owner_share;

-- +goose StatementEnd
