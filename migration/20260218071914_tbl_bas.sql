-- +goose Up
-- +goose StatementBegin

-- ============================================================
-- BAS Snapshot Header
-- One row per generated Business Activity Statement snapshot.
-- Scoped to a clinic + date range, aligned with ATO BAS labels.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_bas_snapshot (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    clinic_id           UUID NOT NULL
                            REFERENCES tbl_clinic(id) ON DELETE CASCADE,

    quarter_id          UUID
                            REFERENCES tbl_quarter(id) ON DELETE SET NULL,

    -- Period
    period_start        DATE NOT NULL,
    period_end          DATE NOT NULL,
    period_type         VARCHAR(20) NOT NULL DEFAULT 'QUARTERLY'
                            CHECK (period_type IN ('QUARTERLY', 'ANNUALLY')),

    -- ========== GST on Sales (G1–G9) ==========
    g1_total_sales              NUMERIC(14,2) NOT NULL DEFAULT 0,
    g2_export_sales             NUMERIC(14,2) NOT NULL DEFAULT 0,
    g3_gst_free_sales           NUMERIC(14,2) NOT NULL DEFAULT 0,
    g4_input_taxed_sales        NUMERIC(14,2) NOT NULL DEFAULT 0,
    g5_g2_g3_g4                 NUMERIC(14,2) NOT NULL DEFAULT 0,
    g6_total_taxable_sales      NUMERIC(14,2) NOT NULL DEFAULT 0,
    g7_adjustments              NUMERIC(14,2) NOT NULL DEFAULT 0,
    g8_total_taxable_supplies   NUMERIC(14,2) NOT NULL DEFAULT 0,
    g9_gst_on_sales             NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== GST on Purchases (G10–G20) ==========
    g10_capital_purchases       NUMERIC(14,2) NOT NULL DEFAULT 0,
    g11_non_capital_purchases   NUMERIC(14,2) NOT NULL DEFAULT 0,
    g12_g10_g11                 NUMERIC(14,2) NOT NULL DEFAULT 0,
    g13_purchases_for_input_taxed_sales
                                NUMERIC(14,2) NOT NULL DEFAULT 0,
    g14_purchases_without_gst   NUMERIC(14,2) NOT NULL DEFAULT 0,
    g15_private_use             NUMERIC(14,2) NOT NULL DEFAULT 0,
    g16_g13_g14_g15             NUMERIC(14,2) NOT NULL DEFAULT 0,
    g17_total_creditable_purchases
                                NUMERIC(14,2) NOT NULL DEFAULT 0,
    g18_adjustments             NUMERIC(14,2) NOT NULL DEFAULT 0,
    g19_total_creditable_acquisitions
                                NUMERIC(14,2) NOT NULL DEFAULT 0,
    g20_gst_on_purchases        NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== Summary Labels (1A, 1B) ==========
    label_1a_gst_on_sales       NUMERIC(14,2) NOT NULL DEFAULT 0,
    label_1b_gst_on_purchases   NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== PAYG Withholding (W1–W4) ==========
    w1_total_salary_wages       NUMERIC(14,2) NOT NULL DEFAULT 0,
    w2_amounts_withheld         NUMERIC(14,2) NOT NULL DEFAULT 0,
    w3_other_amounts_withheld   NUMERIC(14,2) NOT NULL DEFAULT 0,
    w4_total_withheld           NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== PAYG Instalments (T1–T4) ==========
    t1_instalment_income        NUMERIC(14,2) NOT NULL DEFAULT 0,
    t2_instalment_rate          NUMERIC(7,4)  NOT NULL DEFAULT 0,
    t3_new_varied_rate          NUMERIC(7,4)  NOT NULL DEFAULT 0,
    t4_instalment_amount        NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== Net Amounts ==========
    net_gst_payable             NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount_owing          NUMERIC(14,2) NOT NULL DEFAULT 0,

    -- ========== Lifecycle ==========
    status              VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
                            CHECK (status IN ('DRAFT', 'FINALISED', 'LOCKED')),

    generated_by        UUID NOT NULL
                            REFERENCES tbl_user(id),

    finalised_at        TIMESTAMPTZ,
    finalised_by        UUID
                            REFERENCES tbl_user(id),
    locked_at           TIMESTAMPTZ,
    locked_by           UUID
                            REFERENCES tbl_user(id),

    notes               TEXT,

    snapshot_data       JSONB NOT NULL DEFAULT '{}',

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at          TIMESTAMPTZ
);

-- ============================================================
-- BAS Snapshot Line Items
-- One row per COA account contribution to a BAS label.
-- Provides drill-down from aggregated BAS labels to accounts.
-- ============================================================
CREATE TABLE IF NOT EXISTS tbl_bas_snapshot_line (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    bas_snapshot_id     UUID NOT NULL
                            REFERENCES tbl_bas_snapshot(id) ON DELETE CASCADE,

    coa_id              UUID NOT NULL
                            REFERENCES tbl_account(id) ON DELETE RESTRICT,

    account_code        VARCHAR(10) NOT NULL,
    account_name        VARCHAR(255) NOT NULL,
    account_tax_name    VARCHAR(50) NOT NULL,

    bas_label           VARCHAR(20) NOT NULL
                            CHECK (bas_label IN (
                                'G1','G2','G3','G4','G10','G11','G13','G14',
                                '1A','1B','W1','W2','BAS_EXCLUDED'
                            )),

    base_amount         NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount          NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount        NUMERIC(14,2) NOT NULL DEFAULT 0,

    transaction_count   INT NOT NULL DEFAULT 0,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- Indexes
-- ============================================================

-- tbl_bas_snapshot
CREATE INDEX IF NOT EXISTS idx_bas_snapshot_clinic
    ON tbl_bas_snapshot(clinic_id);

CREATE INDEX IF NOT EXISTS idx_bas_snapshot_clinic_period
    ON tbl_bas_snapshot(clinic_id, period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_bas_snapshot_quarter
    ON tbl_bas_snapshot(quarter_id)
    WHERE quarter_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_bas_snapshot_status
    ON tbl_bas_snapshot(status);

CREATE INDEX IF NOT EXISTS idx_bas_snapshot_active
    ON tbl_bas_snapshot(id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_bas_snapshot_clinic_period_active
    ON tbl_bas_snapshot(clinic_id, period_start, period_end)
    WHERE deleted_at IS NULL;

-- tbl_bas_snapshot_line
CREATE INDEX IF NOT EXISTS idx_bas_line_snapshot
    ON tbl_bas_snapshot_line(bas_snapshot_id);

CREATE INDEX IF NOT EXISTS idx_bas_line_coa
    ON tbl_bas_snapshot_line(coa_id);

CREATE INDEX IF NOT EXISTS idx_bas_line_label
    ON tbl_bas_snapshot_line(bas_label);

CREATE UNIQUE INDEX IF NOT EXISTS ux_bas_line_snapshot_coa_label
    ON tbl_bas_snapshot_line(bas_snapshot_id, coa_id, bas_label);

-- ============================================================
-- Triggers
-- ============================================================
CREATE TRIGGER trg_bas_snapshot_updated_at
    BEFORE UPDATE ON tbl_bas_snapshot
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_bas_snapshot_updated_at ON tbl_bas_snapshot;

DROP INDEX IF EXISTS ux_bas_line_snapshot_coa_label;
DROP INDEX IF EXISTS idx_bas_line_label;
DROP INDEX IF EXISTS idx_bas_line_coa;
DROP INDEX IF EXISTS idx_bas_line_snapshot;

DROP INDEX IF EXISTS ux_bas_snapshot_clinic_period_active;
DROP INDEX IF EXISTS idx_bas_snapshot_active;
DROP INDEX IF EXISTS idx_bas_snapshot_status;
DROP INDEX IF EXISTS idx_bas_snapshot_quarter;
DROP INDEX IF EXISTS idx_bas_snapshot_clinic_period;
DROP INDEX IF EXISTS idx_bas_snapshot_clinic;

DROP TABLE IF EXISTS tbl_bas_snapshot_line;
DROP TABLE IF EXISTS tbl_bas_snapshot;

-- +goose StatementEnd
