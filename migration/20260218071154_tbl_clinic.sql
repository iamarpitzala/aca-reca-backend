-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_clinic (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    user_id VARCHAR(40) REFERENCES tbl_user(id),

    name VARCHAR(255) NOT NULL,
    abn_number VARCHAR(11) NOT NULL,

    address TEXT NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,

    postcode VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),

    logo_url TEXT,
    description TEXT,

    share_type VARCHAR(50) NOT NULL DEFAULT 'PERCENTAGE',
    clinic_share INT NOT NULL DEFAULT 50,
    owner_share INT NOT NULL DEFAULT 50,
    method_type VARCHAR(50) NOT NULL DEFAULT 'NET' CHECK (method_type IN ('NET', 'GROSS')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT chk_share_type CHECK (share_type IN ('FIXED', 'PERCENTAGE')),
    CONSTRAINT chk_share_percentage CHECK (
        share_type != 'PERCENTAGE'
        OR (clinic_share + owner_share = 100)
    )
);


CREATE UNIQUE INDEX ux_clinic_abn_active
ON tbl_clinic (abn_number)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_clinic_abn_active;
DROP TABLE IF EXISTS tbl_clinic;
-- +goose StatementEnd
