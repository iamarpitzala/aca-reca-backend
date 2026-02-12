-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    ADD COLUMN method_type VARCHAR(50) NOT NULL DEFAULT 'NET' CHECK (method_type IN ('NET', 'GROSS'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_clinic
    DROP COLUMN IF EXISTS method_type;
-- +goose StatementEnd
