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

func (r *basSnapshotRepo) Create(ctx context.Context, snapshot *domain.BASSnapshot) error {
	snapshot.ID = uuid.New()
	snapshot.CreatedAt = time.Now()
	snapshot.UpdatedAt = time.Now()
	
	query := `INSERT INTO tbl_bas_snapshot 
		(id, clinic_id, period_start, period_end, period_type, g1_total_sales, g2_export_sales, g3_gst_free_sales, g10_capital_purchases, g11_non_capital_purchases, label_1a_gst_on_sales, label_1b_gst_on_purchases, net_gst_payable, status, finalised_at, finalised_by, snapshot_data, created_at, updated_at)
		VALUES (:id, :clinic_id, :period_start, :period_end, :period_type, :g1_total_sales, :g2_export_sales, :g3_gst_free_sales, :g10_capital_purchases, :g11_non_capital_purchases, :label_1a_gst_on_sales, :label_1b_gst_on_purchases, :net_gst_payable, :status, :finalised_at, :finalised_by, :snapshot_data, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, snapshot)
	return err
}

func (r *basSnapshotRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error) {
	query := `SELECT id, clinic_id, period_start, period_end, period_type, g1_total_sales, g2_export_sales, g3_gst_free_sales, g10_capital_purchases, g11_non_capital_purchases, label_1a_gst_on_sales, label_1b_gst_on_purchases, net_gst_payable, status, finalised_at, finalised_by, snapshot_data, created_at, updated_at, deleted_at
		FROM tbl_bas_snapshot 
		WHERE id = $1 AND deleted_at IS NULL`
	var snapshot domain.BASSnapshot
	err := r.db.GetContext(ctx, &snapshot, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("BAS snapshot not found")
		}
		return nil, errors.New("failed to get BAS snapshot")
	}
	return &snapshot, nil
}

func (r *basSnapshotRepo) GetByClinicIDAndPeriod(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) (*domain.BASSnapshot, error) {
	query := `SELECT id, clinic_id, period_start, period_end, period_type, g1_total_sales, g2_export_sales, g3_gst_free_sales, g10_capital_purchases, g11_non_capital_purchases, label_1a_gst_on_sales, label_1b_gst_on_purchases, net_gst_payable, status, finalised_at, finalised_by, snapshot_data, created_at, updated_at, deleted_at
		FROM tbl_bas_snapshot 
		WHERE clinic_id = $1 AND period_start = $2 AND period_end = $3 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1`
	var snapshot domain.BASSnapshot
	err := r.db.GetContext(ctx, &snapshot, query, clinicID, periodStart, periodEnd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("BAS snapshot not found")
		}
		return nil, errors.New("failed to get BAS snapshot")
	}
	return &snapshot, nil
}

func (r *basSnapshotRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error) {
	query := `SELECT id, clinic_id, period_start, period_end, period_type, g1_total_sales, g2_export_sales, g3_gst_free_sales, g10_capital_purchases, g11_non_capital_purchases, label_1a_gst_on_sales, label_1b_gst_on_purchases, net_gst_payable, status, finalised_at, finalised_by, snapshot_data, created_at, updated_at, deleted_at
		FROM tbl_bas_snapshot 
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY period_start DESC`
	var snapshots []domain.BASSnapshot
	err := r.db.SelectContext(ctx, &snapshots, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to get BAS snapshots")
	}
	return snapshots, nil
}

func (r *basSnapshotRepo) GetFinalisedByClinicIDs(ctx context.Context, clinicIDs []uuid.UUID, periodStart, periodEnd time.Time) ([]domain.BASSnapshot, error) {
	if len(clinicIDs) == 0 {
		return []domain.BASSnapshot{}, nil
	}
	
	query, args, err := sqlx.In(`SELECT id, clinic_id, period_start, period_end, period_type, g1_total_sales, g2_export_sales, g3_gst_free_sales, g10_capital_purchases, g11_non_capital_purchases, label_1a_gst_on_sales, label_1b_gst_on_purchases, net_gst_payable, status, finalised_at, finalised_by, snapshot_data, created_at, updated_at, deleted_at
		FROM tbl_bas_snapshot 
		WHERE clinic_id IN (?) AND period_start = ? AND period_end = ? AND status IN ('FINALISED', 'LOCKED') AND deleted_at IS NULL
		ORDER BY clinic_id, period_start DESC`, clinicIDs, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	
	var snapshots []domain.BASSnapshot
	err = r.db.SelectContext(ctx, &snapshots, query, args...)
	if err != nil {
		return nil, errors.New("failed to get finalised BAS snapshots")
	}
	return snapshots, nil
}

func (r *basSnapshotRepo) Update(ctx context.Context, snapshot *domain.BASSnapshot) error {
	snapshot.UpdatedAt = time.Now()
	query := `UPDATE tbl_bas_snapshot SET
		g1_total_sales = :g1_total_sales,
		g2_export_sales = :g2_export_sales,
		g3_gst_free_sales = :g3_gst_free_sales,
		g10_capital_purchases = :g10_capital_purchases,
		g11_non_capital_purchases = :g11_non_capital_purchases,
		label_1a_gst_on_sales = :label_1a_gst_on_sales,
		label_1b_gst_on_purchases = :label_1b_gst_on_purchases,
		net_gst_payable = :net_gst_payable,
		status = :status,
		finalised_at = :finalised_at,
		finalised_by = :finalised_by,
		snapshot_data = :snapshot_data,
		updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, snapshot)
	return err
}

func (r *basSnapshotRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_bas_snapshot SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
