-- +goose Up
-- +goose StatementBegin
ALTER TABLE tbl_clinic ADD COLUMN user_id VARCHAR(40) REFERENCES tbl_user(id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_clinic DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd
