package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type basSnapshotRepo struct {
	db *sqlx.DB
}

func NewBASSnapshotRepository(db *sqlx.DB) port.BASSnapshotRepository {
	return &basSnapshotRepo{db: db}
}

func (r *basSnapshotRepo) Create(ctx context.Context, s *domain.BASSnapshot) error {
	q := `
		INSERT INTO tbl_bas_snapshot (
			id, clinic_id, quarter_id,
			period_start, period_end, period_type,
			g1_total_sales, g2_export_sales, g3_gst_free_sales,
			g4_input_taxed_sales, g5_g2_g3_g4,
			g6_total_taxable_sales, g7_adjustments,
			g8_total_taxable_supplies, g9_gst_on_sales,
			g10_capital_purchases, g11_non_capital_purchases, g12_g10_g11,
			g13_purchases_for_input_taxed_sales, g14_purchases_without_gst,
			g15_private_use, g16_g13_g14_g15,
			g17_total_creditable_purchases, g18_adjustments,
			g19_total_creditable_acquisitions, g20_gst_on_purchases,
			label_1a_gst_on_sales, label_1b_gst_on_purchases,
			w1_total_salary_wages, w2_amounts_withheld,
			w3_other_amounts_withheld, w4_total_withheld,
			t1_instalment_income, t2_instalment_rate,
			t3_new_varied_rate, t4_instalment_amount,
			net_gst_payable, total_amount_owing,
			status, generated_by,
			finalised_at, finalised_by,
			locked_at, locked_by,
			notes, snapshot_data,
			created_at, updated_at
		) VALUES (
			:id, :clinic_id, :quarter_id,
			:period_start, :period_end, :period_type,
			:g1_total_sales, :g2_export_sales, :g3_gst_free_sales,
			:g4_input_taxed_sales, :g5_g2_g3_g4,
			:g6_total_taxable_sales, :g7_adjustments,
			:g8_total_taxable_supplies, :g9_gst_on_sales,
			:g10_capital_purchases, :g11_non_capital_purchases, :g12_g10_g11,
			:g13_purchases_for_input_taxed_sales, :g14_purchases_without_gst,
			:g15_private_use, :g16_g13_g14_g15,
			:g17_total_creditable_purchases, :g18_adjustments,
			:g19_total_creditable_acquisitions, :g20_gst_on_purchases,
			:label_1a_gst_on_sales, :label_1b_gst_on_purchases,
			:w1_total_salary_wages, :w2_amounts_withheld,
			:w3_other_amounts_withheld, :w4_total_withheld,
			:t1_instalment_income, :t2_instalment_rate,
			:t3_new_varied_rate, :t4_instalment_amount,
			:net_gst_payable, :total_amount_owing,
			:status, :generated_by,
			:finalised_at, :finalised_by,
			:locked_at, :locked_by,
			:notes, :snapshot_data,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, s)
	return err
}

func (r *basSnapshotRepo) CreateLines(ctx context.Context, lines []domain.BASSnapshotLine) error {
	if len(lines) == 0 {
		return nil
	}
	q := `
		INSERT INTO tbl_bas_snapshot_line (
			id, bas_snapshot_id, coa_id,
			account_code, account_name, account_tax_name,
			bas_label,
			base_amount, gst_amount, total_amount,
			transaction_count, created_at
		) VALUES (
			:id, :bas_snapshot_id, :coa_id,
			:account_code, :account_name, :account_tax_name,
			:bas_label,
			:base_amount, :gst_amount, :total_amount,
			:transaction_count, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, lines)
	return err
}

var basSnapshotColumns = `
	id, clinic_id, quarter_id,
	period_start, period_end, period_type,
	g1_total_sales, g2_export_sales, g3_gst_free_sales,
	g4_input_taxed_sales, g5_g2_g3_g4,
	g6_total_taxable_sales, g7_adjustments,
	g8_total_taxable_supplies, g9_gst_on_sales,
	g10_capital_purchases, g11_non_capital_purchases, g12_g10_g11,
	g13_purchases_for_input_taxed_sales, g14_purchases_without_gst,
	g15_private_use, g16_g13_g14_g15,
	g17_total_creditable_purchases, g18_adjustments,
	g19_total_creditable_acquisitions, g20_gst_on_purchases,
	label_1a_gst_on_sales, label_1b_gst_on_purchases,
	w1_total_salary_wages, w2_amounts_withheld,
	w3_other_amounts_withheld, w4_total_withheld,
	t1_instalment_income, t2_instalment_rate,
	t3_new_varied_rate, t4_instalment_amount,
	net_gst_payable, total_amount_owing,
	status, generated_by,
	finalised_at, finalised_by,
	locked_at, locked_by,
	notes, snapshot_data,
	created_at, updated_at, deleted_at
`

func (r *basSnapshotRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error) {
	q := `SELECT ` + basSnapshotColumns + `
		FROM tbl_bas_snapshot
		WHERE id = $1 AND deleted_at IS NULL`

	var snapshot domain.BASSnapshot
	if err := r.db.GetContext(ctx, &snapshot, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get BAS snapshot")
	}
	return &snapshot, nil
}

func (r *basSnapshotRepo) GetLinesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) ([]domain.BASSnapshotLine, error) {
	q := `
		SELECT id, bas_snapshot_id, coa_id,
		       account_code, account_name, account_tax_name,
		       bas_label,
		       base_amount, gst_amount, total_amount,
		       transaction_count, created_at
		FROM tbl_bas_snapshot_line
		WHERE bas_snapshot_id = $1
		ORDER BY bas_label, account_code
	`
	var lines []domain.BASSnapshotLine
	if err := r.db.SelectContext(ctx, &lines, q, snapshotID); err != nil {
		return nil, errors.New("failed to get BAS snapshot lines")
	}
	return lines, nil
}

func (r *basSnapshotRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error) {
	q := `SELECT ` + basSnapshotColumns + `
		FROM tbl_bas_snapshot
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY period_start DESC, created_at DESC`

	var snapshots []domain.BASSnapshot
	if err := r.db.SelectContext(ctx, &snapshots, q, clinicID); err != nil {
		return nil, errors.New("failed to list BAS snapshots for clinic")
	}
	return snapshots, nil
}

func (r *basSnapshotRepo) Update(ctx context.Context, s *domain.BASSnapshot) error {
	q := `
		UPDATE tbl_bas_snapshot SET
			period_start = :period_start,
			period_end = :period_end,
			period_type = :period_type,
			g1_total_sales = :g1_total_sales,
			g2_export_sales = :g2_export_sales,
			g3_gst_free_sales = :g3_gst_free_sales,
			g4_input_taxed_sales = :g4_input_taxed_sales,
			g5_g2_g3_g4 = :g5_g2_g3_g4,
			g6_total_taxable_sales = :g6_total_taxable_sales,
			g7_adjustments = :g7_adjustments,
			g8_total_taxable_supplies = :g8_total_taxable_supplies,
			g9_gst_on_sales = :g9_gst_on_sales,
			g10_capital_purchases = :g10_capital_purchases,
			g11_non_capital_purchases = :g11_non_capital_purchases,
			g12_g10_g11 = :g12_g10_g11,
			g13_purchases_for_input_taxed_sales = :g13_purchases_for_input_taxed_sales,
			g14_purchases_without_gst = :g14_purchases_without_gst,
			g15_private_use = :g15_private_use,
			g16_g13_g14_g15 = :g16_g13_g14_g15,
			g17_total_creditable_purchases = :g17_total_creditable_purchases,
			g18_adjustments = :g18_adjustments,
			g19_total_creditable_acquisitions = :g19_total_creditable_acquisitions,
			g20_gst_on_purchases = :g20_gst_on_purchases,
			label_1a_gst_on_sales = :label_1a_gst_on_sales,
			label_1b_gst_on_purchases = :label_1b_gst_on_purchases,
			w1_total_salary_wages = :w1_total_salary_wages,
			w2_amounts_withheld = :w2_amounts_withheld,
			w3_other_amounts_withheld = :w3_other_amounts_withheld,
			w4_total_withheld = :w4_total_withheld,
			t1_instalment_income = :t1_instalment_income,
			t2_instalment_rate = :t2_instalment_rate,
			t3_new_varied_rate = :t3_new_varied_rate,
			t4_instalment_amount = :t4_instalment_amount,
			net_gst_payable = :net_gst_payable,
			total_amount_owing = :total_amount_owing,
			notes = :notes,
			snapshot_data = :snapshot_data,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	res, err := r.db.NamedExecContext(ctx, q, s)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("BAS snapshot not found")
	}
	return nil
}

func (r *basSnapshotRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, finalisedAt *time.Time, finalisedBy *uuid.UUID) error {
	q := `
		UPDATE tbl_bas_snapshot SET
			status = $1,
			finalised_at = $2,
			finalised_by = $3,
			updated_at = now()
		WHERE id = $4 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, status, finalisedAt, finalisedBy, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("BAS snapshot not found")
	}
	return nil
}

func (r *basSnapshotRepo) UpdateLock(ctx context.Context, id uuid.UUID, lockedAt *time.Time, lockedBy *uuid.UUID) error {
	q := `
		UPDATE tbl_bas_snapshot SET
			status = 'LOCKED',
			locked_at = $1,
			locked_by = $2,
			updated_at = now()
		WHERE id = $3 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, lockedAt, lockedBy, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("BAS snapshot not found")
	}
	return nil
}

func (r *basSnapshotRepo) DeleteLinesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM tbl_bas_snapshot_line WHERE bas_snapshot_id = $1`,
		snapshotID,
	)
	return err
}

func (r *basSnapshotRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
		UPDATE tbl_bas_snapshot SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("BAS snapshot not found")
	}
	return nil
}

func (r *basSnapshotRepo) AggregateLedger(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) ([]domain.BASLedgerAggRow, error) {
	q := `
		SELECT
			tl.coa_id,
			a.code           AS account_code,
			a.name           AS account_name,
			at_type.name     AS account_type,
			at_tax.name      AS account_tax_name,
			COALESCE(SUM(CASE WHEN tl.entry_type = 'DEBIT'  THEN tl.amount ELSE 0 END), 0) AS debit_total,
			COALESCE(SUM(CASE WHEN tl.entry_type = 'CREDIT' THEN tl.amount ELSE 0 END), 0) AS credit_total,
			COALESCE(SUM(tl.net_amount), 0)  AS net_amount,
			COALESCE(SUM(tl.gst_amount), 0)  AS gst_amount,
			COUNT(DISTINCT tl.transaction_id) AS transaction_count
		FROM tbl_transaction_ledger tl
		JOIN tbl_transaction t         ON t.id = tl.transaction_id
		JOIN tbl_account a             ON a.id = tl.coa_id
		JOIN tbl_account_type at_type  ON at_type.id = a.account_type_id
		JOIN tbl_account_tax at_tax    ON at_tax.id = a.account_tax_id
		WHERE t.clinic_id   = $1
		  AND t.status       = 'POSTED'
		  AND t.deleted_at   IS NULL
		  AND tl.transaction_date BETWEEN $2 AND $3
		  AND at_tax.name   != 'BAS Excluded'
		GROUP BY tl.coa_id, a.code, a.name, at_type.name, at_tax.name
		ORDER BY a.code
	`
	var rows []domain.BASLedgerAggRow
	if err := r.db.SelectContext(ctx, &rows, q, clinicID, periodStart, periodEnd); err != nil {
		return nil, errors.New("failed to aggregate ledger for BAS")
	}
	return rows, nil
}
