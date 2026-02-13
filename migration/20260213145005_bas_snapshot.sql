-- +goose Up
-- +goose StatementBegin

-- BAS Snapshot table for finalised BAS records
CREATE TABLE IF NOT EXISTS tbl_bas_snapshot (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    
    -- Period identification
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('QUARTERLY', 'ANNUALLY')),
    
    -- BAS Fields (ATO standard)
    g1_total_sales NUMERIC(14,2) NOT NULL DEFAULT 0,
    g2_export_sales NUMERIC(14,2) NOT NULL DEFAULT 0,
    g3_gst_free_sales NUMERIC(14,2) NOT NULL DEFAULT 0,
    g10_capital_purchases NUMERIC(14,2) NOT NULL DEFAULT 0,
    g11_non_capital_purchases NUMERIC(14,2) NOT NULL DEFAULT 0,
    label_1a_gst_on_sales NUMERIC(14,2) NOT NULL DEFAULT 0,
    label_1b_gst_on_purchases NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_gst_payable NUMERIC(14,2) NOT NULL DEFAULT 0,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'FINALISED', 'LOCKED')),
    
    -- Finalisation metadata
    finalised_at TIMESTAMP NULL,
    finalised_by UUID NULL REFERENCES tbl_user(id) ON DELETE SET NULL,
    
    -- Snapshot data (full JSONB snapshot of all related transactions/entries)
    snapshot_data JSONB NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Unique constraint: one finalised BAS per clinic per period
CREATE UNIQUE INDEX ux_bas_snapshot_clinic_period_finalised
ON tbl_bas_snapshot (clinic_id, period_start, period_end)
WHERE status IN ('FINALISED', 'LOCKED') AND deleted_at IS NULL;

-- Indexes for lookups
CREATE INDEX idx_bas_snapshot_clinic_id ON tbl_bas_snapshot(clinic_id);
CREATE INDEX idx_bas_snapshot_period ON tbl_bas_snapshot(period_start, period_end);
CREATE INDEX idx_bas_snapshot_status ON tbl_bas_snapshot(status);
CREATE INDEX idx_bas_snapshot_deleted_at ON tbl_bas_snapshot(deleted_at);

-- Trigger for updated_at
CREATE TRIGGER trg_bas_snapshot_updated_at
BEFORE UPDATE ON tbl_bas_snapshot
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_bas_snapshot_updated_at ON tbl_bas_snapshot;
DROP INDEX IF EXISTS idx_bas_snapshot_deleted_at;
DROP INDEX IF EXISTS idx_bas_snapshot_status;
DROP INDEX IF EXISTS idx_bas_snapshot_period;
DROP INDEX IF EXISTS idx_bas_snapshot_clinic_id;
DROP INDEX IF EXISTS ux_bas_snapshot_clinic_period_finalised;
DROP TABLE IF EXISTS tbl_bas_snapshot;

-- +goose StatementEnd
