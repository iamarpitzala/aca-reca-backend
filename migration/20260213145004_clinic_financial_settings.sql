-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_clinic_financial_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    
    -- Financial Year
    financial_year_start VARCHAR(10) NOT NULL DEFAULT 'JULY' CHECK (financial_year_start IN ('JULY', 'JANUARY')),
    
    -- Accounting Method
    accounting_method VARCHAR(20) NOT NULL DEFAULT 'CASH' CHECK (accounting_method IN ('CASH', 'ACCRUAL')),
    
    -- GST Registration
    gst_registered BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- GST Reporting Frequency
    gst_reporting_frequency VARCHAR(20) NOT NULL DEFAULT 'QUARTERLY' CHECK (gst_reporting_frequency IN ('QUARTERLY', 'ANNUALLY')),
    
    -- Default Amount Entry Mode
    default_amount_mode VARCHAR(20) NOT NULL DEFAULT 'GST_INCLUSIVE' CHECK (default_amount_mode IN ('GST_INCLUSIVE', 'GST_EXCLUSIVE')),
    
    -- Lock Date (prevents edits before this date)
    lock_date DATE,
    
    -- GST Defaults (JSONB for flexible mapping)
    -- Example: {"patient_fees": "GST_FREE", "service_fees": "GST_10", "lab_fees": "GST_FREE"}
    gst_defaults JSONB NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Unique constraint: one financial settings per clinic (active only)
CREATE UNIQUE INDEX ux_clinic_financial_settings_clinic_active
ON tbl_clinic_financial_settings (clinic_id)
WHERE deleted_at IS NULL;

-- Index for lookups
CREATE INDEX idx_clinic_financial_settings_clinic_id ON tbl_clinic_financial_settings(clinic_id);
CREATE INDEX idx_clinic_financial_settings_deleted_at ON tbl_clinic_financial_settings(deleted_at);

-- Trigger for updated_at
CREATE TRIGGER trg_clinic_financial_settings_updated_at
BEFORE UPDATE ON tbl_clinic_financial_settings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_clinic_financial_settings_updated_at ON tbl_clinic_financial_settings;
DROP INDEX IF EXISTS idx_clinic_financial_settings_deleted_at;
DROP INDEX IF EXISTS idx_clinic_financial_settings_clinic_id;
DROP INDEX IF EXISTS ux_clinic_financial_settings_clinic_active;
DROP TABLE IF EXISTS tbl_clinic_financial_settings;

-- +goose StatementEnd
