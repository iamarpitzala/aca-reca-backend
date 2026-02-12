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

type clinicRepo struct {
	db *sqlx.DB
}

func NewClinicRepository(db *sqlx.DB) port.ClinicRepository {
	return &clinicRepo{db: db}
}

func (r *clinicRepo) Create(ctx context.Context, clinic *domain.Clinic) error {
	query := `INSERT INTO tbl_clinic (id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share, method_type, is_active, created_at, updated_at)
		VALUES (:id, :name, :abn_number, :address, :city, :state, :postcode, :phone, :email, :website, :logo_url, :description, :share_type, :clinic_share, :owner_share, :method_type, :is_active, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, clinic)
	return err
}

func (r *clinicRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share, method_type, is_active, created_at, updated_at FROM tbl_clinic WHERE id = $1 AND deleted_at IS NULL`
	var clinic domain.Clinic
	err := r.db.GetContext(ctx, &clinic, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("clinic not found")
		}
		return nil, errors.New("failed to get clinic by id")
	}
	return &clinic, nil
}

func (r *clinicRepo) Update(ctx context.Context, clinic *domain.Clinic) error {
	query := `UPDATE tbl_clinic SET name = :name, abn_number = :abn_number, address = :address, city = :city, state = :state, postcode = :postcode, phone = :phone, email = :email, website = :website, logo_url = :logo_url, description = :description, share_type = :share_type, clinic_share = :clinic_share, owner_share = :owner_share, method_type = :method_type, is_active = :is_active, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, clinic)
	return err
}

func (r *clinicRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_clinic SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *clinicRepo) List(ctx context.Context) ([]domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share, method_type, is_active, created_at, updated_at FROM tbl_clinic WHERE deleted_at IS NULL`
	var clinics []domain.Clinic
	err := r.db.SelectContext(ctx, &clinics, query)
	if err != nil {
		return nil, errors.New("failed to get all clinics")
	}
	return clinics, nil
}

func (r *clinicRepo) GetByABN(ctx context.Context, abnNumber string) (*domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share, method_type, is_active, created_at, updated_at FROM tbl_clinic WHERE abn_number = $1 AND deleted_at IS NULL`
	var clinic domain.Clinic
	err := r.db.GetContext(ctx, &clinic, query, abnNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("clinic not found")
		}
		return nil, errors.New("failed to get clinic by abn number")
	}
	return &clinic, nil
}
