package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type userClinicRepo struct {
	db *sqlx.DB
}

func NewUserClinicRepository(db *sqlx.DB) port.UserClinicRepository {
	return &userClinicRepo{db: db}
}

func (r *userClinicRepo) Create(ctx context.Context, uc *user.UserClinic) error {
	query := `INSERT INTO tbl_user_clinic (id, user_id, clinic_id, role, created_at, updated_at)
		VALUES (:id, :user_id, :clinic_id, :role, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, uc)
	return err
}

func (r *userClinicRepo) GetByID(ctx context.Context, id string) (*user.UserClinic, error) {
	query := `SELECT id, user_id, clinic_id, role, created_at, updated_at FROM tbl_user_clinic WHERE id = $1 AND deleted_at IS NULL`
	var uc user.UserClinic
	err := r.db.GetContext(ctx, &uc, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user clinic not found")
		}
		return nil, errors.New("failed to get user clinic by id")
	}
	return &uc, nil
}

func (r *userClinicRepo) GetByUserAndClinic(ctx context.Context, userID, clinicID string) (*user.UserClinic, error) {
	query := `SELECT id, user_id, clinic_id, role, created_at, updated_at FROM tbl_user_clinic WHERE user_id = $1 AND clinic_id = $2 AND deleted_at IS NULL`
	var uc user.UserClinic
	err := r.db.GetContext(ctx, &uc, query, userID, clinicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get user clinic")
	}
	return &uc, nil
}

func (r *userClinicRepo) GetUserClinics(ctx context.Context, userID string) ([]user.UserClinicWithClinic, error) {
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
	var rows []user.UserClinicWithClinic
	if err := r.db.SelectContext(ctx, &rows, query, userID); err != nil {
		return nil, errors.New("failed to get user clinics")
	}
	out := make([]user.UserClinicWithClinic, 0, len(rows))
	for _, row := range rows {
		out = append(out, user.UserClinicWithClinic{
			UserClinic: user.UserClinic{
				ID:        row.ID,
				UserID:    row.UserID,
				ClinicID:  row.ClinicID,
				Role:      row.Role,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			Clinic: clinic.Clinic{
				ID:          row.ClinicID,
				Name:        row.Name,
				ABNNumber:   row.ABNNumber,
				Address:     row.Address,
				City:        row.City,
				State:       row.State,
				Postcode:    row.Postcode,
				Phone:       row.Phone,
				Email:       row.Email,
				Website:     row.Website,
				LogoURL:     row.LogoURL,
				Description: row.Description,
				ShareType:   row.ShareType,
				ClinicShare: row.ClinicShare,
				OwnerShare:  row.OwnerShare,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			},
		})
	}
	return out, nil
}

func (r *userClinicRepo) GetClinicUsers(ctx context.Context, clinicID string) ([]user.UserClinicWithUser, error) {
	query := `SELECT uc.id, uc.user_id, uc.clinic_id, uc.role, uc.created_at, uc.updated_at,
		u.id as "user.id", u.email as "user.email", u.first_name as "user.firstName",
		u.last_name as "user.lastName", u.phone as "user.phone", u.avatar_url as "user.avatarURL",
		u.is_email_verified as "user.isEmailVerified",
		u.created_at as "user.createdAt", u.updated_at as "user.updatedAt"
		FROM tbl_user_clinic uc
		INNER JOIN tbl_user u ON uc.user_id = u.id
		WHERE uc.clinic_id = $1 AND uc.deleted_at IS NULL AND u.deleted_at IS NULL`
	var out []user.UserClinicWithUser
	if err := r.db.SelectContext(ctx, &out, query, clinicID); err != nil {
		return nil, errors.New("failed to get clinic users")
	}
	return out, nil
}

func (r *userClinicRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE tbl_user_clinic SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
