-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_clinic (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

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

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- -------------------------
-- UNIQUE (ACTIVE CLINICS ONLY)
-- -------------------------
CREATE UNIQUE INDEX ux_clinic_abn_active
ON tbl_clinic (abn_number)
WHERE deleted_at IS NULL;

-- -------------------------
-- COMMON LOOKUPS
-- -------------------------
CREATE INDEX idx_clinic_active
ON tbl_clinic (id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_clinic_name_active
ON tbl_clinic (LOWER(name))
WHERE deleted_at IS NULL;

CREATE INDEX idx_clinic_city_state_active
ON tbl_clinic (city, state)
WHERE deleted_at IS NULL;

CREATE TRIGGER trg_clinic_updated_at
BEFORE UPDATE ON tbl_clinic
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_clinic;
-- +goose StatementEnd
