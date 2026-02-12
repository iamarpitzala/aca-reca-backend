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

type clinicCOARepo struct {
	db *sqlx.DB
}

func NewClinicCOARepository(db *sqlx.DB) port.ClinicCOARepository {
	return &clinicCOARepo{db: db}
}

func (r *clinicCOARepo) Create(ctx context.Context, cc *domain.ClinicCOA) error {
	query := `INSERT INTO tbl_clinic_coa (id, clinic_id, coa_id, created_at, updated_at)
		VALUES (:id, :clinic_id, :coa_id, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, cc)
	return err
}

func (r *clinicCOARepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ClinicCOA, error) {
	query := `SELECT id, clinic_id, coa_id, created_at, updated_at, deleted_at
		FROM tbl_clinic_coa WHERE id = $1 AND deleted_at IS NULL`
	var cc domain.ClinicCOA
	err := r.db.GetContext(ctx, &cc, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}

func (r *clinicCOARepo) ListByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ClinicCOA, error) {
	query := `SELECT id, clinic_id, coa_id, created_at, updated_at, deleted_at
		FROM tbl_clinic_coa WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY created_at`
	var list []domain.ClinicCOA
	err := r.db.SelectContext(ctx, &list, query, clinicID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []domain.ClinicCOA{}
	}
	return list, nil
}

func (r *clinicCOARepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_clinic_coa SET deleted_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *clinicCOARepo) Exists(ctx context.Context, clinicID, coaID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM tbl_clinic_coa WHERE clinic_id = $1 AND coa_id = $2 AND deleted_at IS NULL LIMIT 1`
	var exists int
	err := r.db.GetContext(ctx, &exists, query, clinicID, coaID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("failed to check clinic AOC existence")
	}
	return true, nil
}
