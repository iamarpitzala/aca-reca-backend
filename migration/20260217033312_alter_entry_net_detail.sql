-- +goose Up
-- +goose StatementBegin

ALTER TABLE tbl_entry_net_details
ADD COLUMN IF NOT EXISTS net_amount NUMERIC(14,2) NOT NULL DEFAULT 0 AFTER source_entry_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tbl_entry_net_details DROP COLUMN IF EXISTS net_amount;

-- +goose StatementEnd
