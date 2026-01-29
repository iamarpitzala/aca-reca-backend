-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_gst (
    id SERIAL PRIMARY KEY,  -- auto-increment integer

    name TEXT NOT NULL,

    type TEXT NOT NULL CHECK (type IN ('INCLUSIVE', 'EXCLUSIVE', 'MANUAL')),

    percentage NUMERIC(5,2) NOT NULL,  -- GST value

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    delete_at TIMESTAMP NULL
);

INSERT INTO tbl_gst (name, type, percentage) VALUES
('GST Inclusive', 'INCLUSIVE', 10.00),
('GST Exclusive', 'EXCLUSIVE', 10.00),
('GST Manual',    'MANUAL',    0.00);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tbl_gst
WHERE name IN ('GST Inclusive', 'GST Exclusive', 'GST Manual');

-- +goose StatementEnd
