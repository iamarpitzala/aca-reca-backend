-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_user (
    id VARCHAR(40) PRIMARY KEY NOT NULL UNIQUE,

    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(50),
    avatar_url TEXT,

    is_email_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    deleted_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_user;
-- +goose StatementEnd
