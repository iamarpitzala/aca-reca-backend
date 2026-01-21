-- +goose Up
-- +goose StatementBegin
create table tbl_clinic (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    abn_number VARCHAR(11) NOT NULL,

    address TEXT NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,

    postcode VARCHAR(255) NULL,
    phone VARCHAR(255) NULL,
    email VARCHAR(255) NULL,
    website VARCHAR(255) NULL,  

    logo_url TEXT NULL,
    description TEXT NULL,

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_clinic_name ON tbl_clinic(name);
CREATE INDEX idx_clinic_abn_number ON tbl_clinic(abn_number);
CREATE INDEX idx_clinic_is_active ON tbl_clinic(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tbl_clinic;
-- +goose StatementEnd
