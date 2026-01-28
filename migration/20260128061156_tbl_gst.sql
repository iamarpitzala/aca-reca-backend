-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_gst (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  -- auto-increment integer

    name TEXT NOT NULL,

    type TEXT NOT NULL
        CHECK (type IN ('INCLUSIVE', 'EXCLUSIVE', 'MANUAL')),

    percentage NUMERIC(5,2) NOT NULL,  -- GST value

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Seed GST rows
INSERT INTO tbl_gst (name, type, percentage) VALUES
('GST Inclusive', 'INCLUSIVE', 10.00),
('GST Exclusive', 'EXCLUSIVE', 10.00),
('GST Manual',    'MANUAL',    0.00);

ALTER TABLE tbl_financial_form
ADD COLUMN gst_id INT DEFAULT NULL;

ALTER TABLE tbl_financial_form
ADD CONSTRAINT fk_tbl_financial_form_gst
FOREIGN KEY (gst_id)
REFERENCES tbl_gst(id)
ON DELETE CASCADE;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete seeded rows first (optional, in case other tables reference them)
DELETE FROM tbl_gst
WHERE name IN ('GST Inclusive', 'GST Exclusive', 'GST Manual');

ALTER TABLE tbl_financial_form
DROP CONSTRAINT IF EXISTS fk_tbl_financial_form_gst;

ALTER TABLE tbl_financial_form
DROP COLUMN IF EXISTS gst_id;
-- +goose StatementEnd

-- +goose StatementBegin
-- Drop table
DROP TABLE IF EXISTS tbl_gst;
-- +goose StatementEnd
