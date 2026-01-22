-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_user_clinic (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'owner',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_user_clinic_unique ON tbl_user_clinic(user_id, clinic_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_clinic_user_id ON tbl_user_clinic(user_id);
CREATE INDEX idx_user_clinic_clinic_id ON tbl_user_clinic(clinic_id);
CREATE INDEX idx_user_clinic_deleted_at ON tbl_user_clinic(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_user_clinic;
-- +goose StatementEnd
