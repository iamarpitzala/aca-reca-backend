-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_user_clinic (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,
    user_id VARCHAR(40) NOT NULL REFERENCES tbl_user(id),
    clinic_id VARCHAR(40) NOT NULL REFERENCES tbl_clinic(id),
    role VARCHAR(50) DEFAULT 'OWNER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_user_clinic_role
        CHECK (role IN ('OWNER', 'MEMBER', 'VIEWER'))
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_user_clinic;
-- +goose StatementEnd
