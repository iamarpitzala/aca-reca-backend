-- +goose Up
-- +goose StatementBegin

-- ============================================================
-- P&L Report Header
-- One row per generated Profit & Loss report.
-- Scoped to a clinic + date range, optionally linked to a quarter.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_pnl_report (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    clinic_id       UUID NOT NULL
                        REFERENCES tbl_clinic(id) ON DELETE CASCADE,

    quarter_id      UUID
                        REFERENCES tbl_quarter(id) ON DELETE SET NULL,

    report_name     VARCHAR(255) NOT NULL,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,

    -- Aggregated totals (denormalised for fast reads)
    total_revenue       NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_cogs          NUMERIC(14,2) NOT NULL DEFAULT 0,
    gross_profit        NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_expenses      NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_profit          NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- GST summary
    total_gst_collected NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_gst_paid      NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_gst             NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- Report lifecycle
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
                        CHECK (status IN ('DRAFT', 'FINAL', 'ARCHIVED')),

    generated_by    UUID NOT NULL
                        REFERENCES tbl_user(id),

    finalized_at    TIMESTAMPTZ,
    notes           TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

-- ============================================================
-- P&L Report Line Items
-- One row per COA account that had activity in the period.
-- Denormalised account_code/name for snapshot immutability.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_pnl_report_line (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    pnl_report_id   UUID NOT NULL
                        REFERENCES tbl_pnl_report(id) ON DELETE CASCADE,

    coa_id          UUID NOT NULL
                        REFERENCES tbl_account(id) ON DELETE RESTRICT,

    account_code    VARCHAR(10) NOT NULL,
    account_name    VARCHAR(255) NOT NULL,

    line_category   VARCHAR(20) NOT NULL
                        CHECK (line_category IN ('REVENUE', 'COGS', 'EXPENSE')),

    debit_total     NUMERIC(14,2) NOT NULL DEFAULT 0,
    credit_total    NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount      NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount      NUMERIC(14,2) NOT NULL DEFAULT 0,

    transaction_count INT NOT NULL DEFAULT 0,

    display_order   INT NOT NULL DEFAULT 0,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- Indexes
-- ============================================================

-- tbl_pnl_report
CREATE INDEX IF NOT EXISTS idx_pnl_report_clinic
    ON tbl_pnl_report(clinic_id);

CREATE INDEX IF NOT EXISTS idx_pnl_report_clinic_period
    ON tbl_pnl_report(clinic_id, period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_pnl_report_quarter
    ON tbl_pnl_report(quarter_id)
    WHERE quarter_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_pnl_report_active
    ON tbl_pnl_report(id)
    WHERE deleted_at IS NULL;

-- tbl_pnl_report_line
CREATE INDEX IF NOT EXISTS idx_pnl_line_report
    ON tbl_pnl_report_line(pnl_report_id);

CREATE INDEX IF NOT EXISTS idx_pnl_line_coa
    ON tbl_pnl_report_line(coa_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_pnl_line_report_coa
    ON tbl_pnl_report_line(pnl_report_id, coa_id);

-- ============================================================
-- Triggers
-- ============================================================
CREATE TRIGGER trg_pnl_report_updated_at
    BEFORE UPDATE ON tbl_pnl_report
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_pnl_report_updated_at ON tbl_pnl_report;

DROP INDEX IF EXISTS ux_pnl_line_report_coa;
DROP INDEX IF EXISTS idx_pnl_line_coa;
DROP INDEX IF EXISTS idx_pnl_line_report;

DROP INDEX IF EXISTS idx_pnl_report_active;
DROP INDEX IF EXISTS idx_pnl_report_quarter;
DROP INDEX IF EXISTS idx_pnl_report_clinic_period;
DROP INDEX IF EXISTS idx_pnl_report_clinic;

DROP TABLE IF EXISTS tbl_pnl_report_line;
DROP TABLE IF EXISTS tbl_pnl_report;

-- +goose StatementEnd
