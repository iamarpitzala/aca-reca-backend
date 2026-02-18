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

type pnlReportRepo struct {
	db *sqlx.DB
}

func NewPnlReportRepository(db *sqlx.DB) port.PnlReportRepository {
	return &pnlReportRepo{db: db}
}

func (r *pnlReportRepo) Create(ctx context.Context, report *domain.PnlReport) error {
	q := `
		INSERT INTO tbl_pnl_report (
			id, clinic_id, quarter_id, report_name,
			period_start, period_end,
			total_revenue, total_cogs, gross_profit,
			total_expenses, net_profit,
			total_gst_collected, total_gst_paid, net_gst,
			status, generated_by, finalized_at, notes,
			created_at, updated_at
		) VALUES (
			:id, :clinic_id, :quarter_id, :report_name,
			:period_start, :period_end,
			:total_revenue, :total_cogs, :gross_profit,
			:total_expenses, :net_profit,
			:total_gst_collected, :total_gst_paid, :net_gst,
			:status, :generated_by, :finalized_at, :notes,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, report)
	return err
}

func (r *pnlReportRepo) CreateLines(ctx context.Context, lines []domain.PnlReportLine) error {
	if len(lines) == 0 {
		return nil
	}
	q := `
		INSERT INTO tbl_pnl_report_line (
			id, pnl_report_id, coa_id,
			account_code, account_name, line_category,
			debit_total, credit_total, net_amount, gst_amount,
			transaction_count, display_order, created_at
		) VALUES (
			:id, :pnl_report_id, :coa_id,
			:account_code, :account_name, :line_category,
			:debit_total, :credit_total, :net_amount, :gst_amount,
			:transaction_count, :display_order, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, lines)
	return err
}

func (r *pnlReportRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PnlReport, error) {
	q := `
		SELECT id, clinic_id, quarter_id, report_name,
		       period_start, period_end,
		       total_revenue, total_cogs, gross_profit,
		       total_expenses, net_profit,
		       total_gst_collected, total_gst_paid, net_gst,
		       status, generated_by, finalized_at, notes,
		       created_at, updated_at, deleted_at
		FROM tbl_pnl_report
		WHERE id = $1 AND deleted_at IS NULL
	`
	var report domain.PnlReport
	if err := r.db.GetContext(ctx, &report, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get P&L report")
	}
	return &report, nil
}

func (r *pnlReportRepo) GetLinesByReportID(ctx context.Context, reportID uuid.UUID) ([]domain.PnlReportLine, error) {
	q := `
		SELECT id, pnl_report_id, coa_id,
		       account_code, account_name, line_category,
		       debit_total, credit_total, net_amount, gst_amount,
		       transaction_count, display_order, created_at
		FROM tbl_pnl_report_line
		WHERE pnl_report_id = $1
		ORDER BY display_order, account_code
	`
	var lines []domain.PnlReportLine
	if err := r.db.SelectContext(ctx, &lines, q, reportID); err != nil {
		return nil, errors.New("failed to get P&L report lines")
	}
	return lines, nil
}

func (r *pnlReportRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.PnlReport, error) {
	q := `
		SELECT id, clinic_id, quarter_id, report_name,
		       period_start, period_end,
		       total_revenue, total_cogs, gross_profit,
		       total_expenses, net_profit,
		       total_gst_collected, total_gst_paid, net_gst,
		       status, generated_by, finalized_at, notes,
		       created_at, updated_at, deleted_at
		FROM tbl_pnl_report
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY period_start DESC, created_at DESC
	`
	var reports []domain.PnlReport
	if err := r.db.SelectContext(ctx, &reports, q, clinicID); err != nil {
		return nil, errors.New("failed to list P&L reports for clinic")
	}
	return reports, nil
}

func (r *pnlReportRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, finalizedAt *time.Time) error {
	q := `
		UPDATE tbl_pnl_report SET
			status = $1,
			finalized_at = $2,
			updated_at = now()
		WHERE id = $3 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, status, finalizedAt, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("P&L report not found")
	}
	return nil
}

func (r *pnlReportRepo) DeleteLinesByReportID(ctx context.Context, reportID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM tbl_pnl_report_line WHERE pnl_report_id = $1`,
		reportID,
	)
	return err
}

func (r *pnlReportRepo) UpdateReport(ctx context.Context, report *domain.PnlReport) error {
	q := `
		UPDATE tbl_pnl_report SET
			total_revenue = :total_revenue,
			total_cogs = :total_cogs,
			gross_profit = :gross_profit,
			total_expenses = :total_expenses,
			net_profit = :net_profit,
			total_gst_collected = :total_gst_collected,
			total_gst_paid = :total_gst_paid,
			net_gst = :net_gst,
			notes = :notes,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	res, err := r.db.NamedExecContext(ctx, q, report)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("P&L report not found")
	}
	return nil
}

func (r *pnlReportRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
		UPDATE tbl_pnl_report SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("P&L report not found")
	}
	return nil
}

// AggregateLedger queries the transaction ledger for POSTED transactions
// in the given date range and aggregates by COA account.
func (r *pnlReportRepo) AggregateLedger(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) ([]domain.LedgerAggRow, error) {
	q := `
		SELECT
			tl.coa_id,
			a.code           AS account_code,
			a.name           AS account_name,
			at.name          AS account_type,
			COALESCE(SUM(CASE WHEN tl.entry_type = 'DEBIT'  THEN tl.amount ELSE 0 END), 0) AS debit_total,
			COALESCE(SUM(CASE WHEN tl.entry_type = 'CREDIT' THEN tl.amount ELSE 0 END), 0) AS credit_total,
			COALESCE(SUM(tl.net_amount), 0)  AS net_amount,
			COALESCE(SUM(tl.gst_amount), 0)  AS gst_amount,
			COUNT(DISTINCT tl.transaction_id) AS transaction_count
		FROM tbl_transaction_ledger tl
		JOIN tbl_transaction t       ON t.id = tl.transaction_id
		JOIN tbl_account a           ON a.id = tl.coa_id
		JOIN tbl_account_type at     ON at.id = a.account_type_id
		WHERE t.clinic_id   = $1
		  AND t.status       = 'POSTED'
		  AND t.deleted_at   IS NULL
		  AND tl.transaction_date BETWEEN $2 AND $3
		  AND at.name IN ('Revenue', 'Expense')
		GROUP BY tl.coa_id, a.code, a.name, at.name
		ORDER BY at.name DESC, a.code
	`
	var rows []domain.LedgerAggRow
	if err := r.db.SelectContext(ctx, &rows, q, clinicID, periodStart, periodEnd); err != nil {
		return nil, errors.New("failed to aggregate ledger for P&L")
	}
	return rows, nil
}
