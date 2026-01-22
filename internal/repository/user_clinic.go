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

func CreateUserClinic(ctx context.Context, db *sqlx.DB, userClinic *domain.UserClinic) error {
	query := `INSERT INTO tbl_user_clinic (id, user_id, clinic_id, role, created_at, updated_at)
		VALUES (:id, :user_id, :clinic_id, :role, :created_at, :updated_at)`
	_, err := db.NamedExecContext(ctx, query, userClinic)
	if err != nil {
		return err
	}
	return nil
}

func GetUserClinicByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.UserClinic, error) {
	query := `SELECT id, user_id, clinic_id, role, created_at, updated_at FROM tbl_user_clinic WHERE id = $1 AND deleted_at IS NULL`
	var userClinic domain.UserClinic
	err := db.GetContext(ctx, &userClinic, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get user clinic by id")
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("user clinic not found")
	}
	return &userClinic, nil
}

func GetUserClinics(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]domain.UserClinicWithClinic, error) {
	query := `SELECT uc.id, uc.user_id, uc.clinic_id, uc.role, uc.created_at, uc.updated_at,
		c.id as "clinic.id", c.name as "clinic.name", c.abn_number as "clinic.abnNumber",
		c.address as "clinic.address", c.city as "clinic.city", c.state as "clinic.state",
		c.postcode as "clinic.postcode", c.phone as "clinic.phone", c.email as "clinic.email",
		c.website as "clinic.website", c.logo_url as "clinic.logoURL", c.description as "clinic.description",
		c.is_active as "clinic.isActive", c.created_at as "clinic.createdAt", c.updated_at as "clinic.updatedAt"
		FROM tbl_user_clinic uc
		INNER JOIN tbl_clinic c ON uc.clinic_id = c.id
		WHERE uc.user_id = $1 AND uc.deleted_at IS NULL AND c.deleted_at IS NULL`

	var userClinics []domain.UserClinicWithClinic
	err := db.SelectContext(ctx, &userClinics, query, userID)
	if err != nil {
		return nil, errors.New("failed to get user clinics")
	}
	return userClinics, nil
}

func GetClinicUsers(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.UserClinicWithUser, error) {
	query := `SELECT uc.id, uc.user_id, uc.clinic_id, uc.role, uc.created_at, uc.updated_at,
		u.id as "user.id", u.email as "user.email", u.first_name as "user.firstName",
		u.last_name as "user.lastName", u.phone as "user.phone", u.avatar_url as "user.avatarURL",
		u.is_email_verified as "user.isEmailVerified",
		u.created_at as "user.createdAt", u.updated_at as "user.updatedAt"
		FROM tbl_user_clinic uc
		INNER JOIN tbl_user u ON uc.user_id = u.id
		WHERE uc.clinic_id = $1 AND uc.deleted_at IS NULL AND u.deleted_at IS NULL`

	var clinicUsers []domain.UserClinicWithUser
	err := db.SelectContext(ctx, &clinicUsers, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to get clinic users")
	}
	return clinicUsers, nil
}

func DeleteUserClinic(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_user_clinic SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func GetUserClinicByUserAndClinic(ctx context.Context, db *sqlx.DB, userID, clinicID uuid.UUID) (*domain.UserClinic, error) {
	query := `SELECT id, user_id, clinic_id, role, created_at, updated_at FROM tbl_user_clinic 
		WHERE user_id = $1 AND clinic_id = $2 AND deleted_at IS NULL`
	var userClinic domain.UserClinic
	err := db.GetContext(ctx, &userClinic, query, userID, clinicID)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get user clinic")
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &userClinic, nil
}
