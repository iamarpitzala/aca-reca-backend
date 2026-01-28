package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateClinic(ctx context.Context, db *sqlx.DB, clinic *domain.Clinic) error {
	query := `INSERT INTO tbl_clinic (id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share,created_at, updated_at)
		VALUES (:id, :name, :abn_number, :address, :city, :state, :postcode, :phone, :email, :website, :logo_url, :description, :share_type, :clinic_share, :owner_share, :created_at, :updated_at)`
	_, err := db.NamedExecContext(ctx, query, clinic)
	if err != nil {
		return err
	}
	return nil
}

func GetClinicByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share,created_at, updated_at FROM tbl_clinic WHERE id = $1 AND deleted_at IS NULL`
	var clinic domain.Clinic
	err := db.GetContext(ctx, &clinic, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get clinic by id")
	}
	return &clinic, nil
}

func UpdateClinic(ctx context.Context, db *sqlx.DB, clinic *domain.Clinic) error {
	query := `UPDATE tbl_clinic SET name = :name, abn_number = :abn_number, address = :address, city = :city, state = :state, postcode = :postcode, phone = :phone, email = :email, website = :website, logo_url = :logo_url, description = :description, share_type = :share_type, clinic_share = :clinic_share, owner_share = :owner_share, updated_at = :updated_at WHERE id = :id`
	_, err := db.NamedExecContext(ctx, query, clinic)
	if err != nil {
		return err
	}
	return nil
}

func DeleteClinic(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_clinic SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllClinics(ctx context.Context, db *sqlx.DB) ([]domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share,created_at, updated_at FROM tbl_clinic WHERE deleted_at IS NULL`
	var clinics []domain.Clinic
	err := db.SelectContext(ctx, &clinics, query)
	if err != nil {
		return nil, errors.New("failed to get all clinics")
	}
	return clinics, nil
}

func GetClinicByABNNumber(ctx context.Context, db *sqlx.DB, abnNumber string) (*domain.Clinic, error) {
	query := `SELECT id, name, abn_number, address, city, state, postcode, phone, email, website, logo_url, description, share_type, clinic_share, owner_share,created_at, updated_at FROM tbl_clinic WHERE abn_number = $1 AND deleted_at IS NULL`
	var clinic domain.Clinic
	err := db.GetContext(ctx, &clinic, query, abnNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get clinic by abn number")
	}
	return &clinic, nil
}
