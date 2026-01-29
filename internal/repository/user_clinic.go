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

// userClinicWithClinicRow is a flat struct for scanning - sqlx doesn't support nested struct scanning
type userClinicWithClinicRow struct {
	UCID         uuid.UUID `db:"uc_id"`
	UCUserID     uuid.UUID `db:"uc_user_id"`
	UCClinicID   uuid.UUID `db:"uc_clinic_id"`
	UCRole       string    `db:"uc_role"`
	UCCreatedAt  time.Time `db:"uc_created_at"`
	UCUpdatedAt  time.Time `db:"uc_updated_at"`
	CID          uuid.UUID `db:"c_id"`
	CName        string    `db:"c_name"`
	CABNNumber   string    `db:"c_abn_number"`
	CAddress     string    `db:"c_address"`
	CCity        string    `db:"c_city"`
	CState       string    `db:"c_state"`
	CPostcode    *string   `db:"c_postcode"`
	CPhone       *string   `db:"c_phone"`
	CEmail       *string   `db:"c_email"`
	CWebsite     *string   `db:"c_website"`
	CLogoURL     *string   `db:"c_logo_url"`
	CDescription *string   `db:"c_description"`
	CShareType   string    `db:"c_share_type"`
	CClinicShare int       `db:"c_clinic_share"`
	COwnerShare  int       `db:"c_owner_share"`
	CCreatedAt   time.Time `db:"c_created_at"`
	CUpdatedAt   time.Time `db:"c_updated_at"`
}

func GetUserClinics(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]domain.UserClinicWithClinic, error) {
	query := `SELECT uc.id as uc_id, uc.user_id as uc_user_id, uc.clinic_id as uc_clinic_id, uc.role as uc_role,
		uc.created_at as uc_created_at, uc.updated_at as uc_updated_at,
		c.id as c_id, c.name as c_name, c.abn_number as c_abn_number, c.address as c_address,
		c.city as c_city, c.state as c_state, c.postcode as c_postcode, c.phone as c_phone,
		c.email as c_email, c.website as c_website, c.logo_url as c_logo_url, c.description as c_description,
		c.share_type as c_share_type, c.clinic_share as c_clinic_share, c.owner_share as c_owner_share,
		c.created_at as c_created_at, c.updated_at as c_updated_at
		FROM tbl_user_clinic uc
		INNER JOIN tbl_clinic c ON uc.clinic_id = c.id
		WHERE uc.user_id = $1 AND uc.deleted_at IS NULL AND c.deleted_at IS NULL`

	var rows []userClinicWithClinicRow
	err := db.SelectContext(ctx, &rows, query, userID)
	if err != nil {
		return nil, errors.New("failed to get user clinics")
	}

	userClinics := make([]domain.UserClinicWithClinic, 0, len(rows))
	for _, r := range rows {
		userClinics = append(userClinics, domain.UserClinicWithClinic{
			UserClinic: domain.UserClinic{
				ID:        r.UCID,
				UserID:    r.UCUserID,
				ClinicID:  r.UCClinicID,
				Role:      r.UCRole,
				CreatedAt: r.UCCreatedAt,
				UpdatedAt: r.UCUpdatedAt,
			},
			Clinic: domain.Clinic{
				ID:          r.CID,
				Name:        r.CName,
				ABNNumber:   r.CABNNumber,
				Address:     r.CAddress,
				City:        r.CCity,
				State:       r.CState,
				Postcode:    r.CPostcode,
				Phone:       r.CPhone,
				Email:       r.CEmail,
				Website:     r.CWebsite,
				LogoURL:     r.CLogoURL,
				Description: r.CDescription,
				ShareType:   r.CShareType,
				ClinicShare: r.CClinicShare,
				OwnerShare:  r.COwnerShare,
				CreatedAt:   r.CCreatedAt,
				UpdatedAt:   r.CUpdatedAt,
			},
		})
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
