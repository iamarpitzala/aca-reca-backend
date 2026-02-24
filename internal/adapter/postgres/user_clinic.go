package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
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
	row := user.UserClinic{
		ID:        uc.ID,
		UserID:    uc.UserID,
		ClinicID:  uc.ClinicID,
		Role:      uc.Role,
		CreatedAt: uc.CreatedAt,
		UpdatedAt: uc.UpdatedAt,
	}
	_, err := r.db.NamedExecContext(ctx, query, &row)
	return err
}

func (r *userClinicRepo) GetByID(ctx context.Context, id string) (*user.UserClinic, error) {
	query := `SELECT id, user_id, clinic_id, role, created_at, updated_at FROM tbl_user_clinic WHERE id = $1 AND deleted_at IS NULL`
	var uc user.UserClinic
	if err := r.db.GetContext(ctx, &uc, query, id); err != nil {
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
	if err := r.db.GetContext(ctx, &uc, query, userID, clinicID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get user clinic")
	}
	return &uc, nil
}

func (r *userClinicRepo) GetUserClinics(ctx context.Context, userID string) ([]user.UserClinicWithClinic, error) {

	query := `
	SELECT
		uc.id as uc_id,
		uc.user_id as uc_user_id,
		uc.clinic_id as uc_clinic_id,
		uc.role as uc_role,
		uc.created_at as uc_created_at,
		uc.updated_at as uc_updated_at,
		c.id as c_id,
		c.name as c_name,
		c.abn_number as c_abn_number,
		c.address as c_address,
		c.city as c_city,
		c.state as c_state,
		c.postcode as c_postcode,
		c.phone as c_phone,
		c.email as c_email,
		c.website as c_website,
		c.logo_url as c_logo_url,
		c.description as c_description,
		c.share_type as c_share_type,
		c.clinic_share as c_clinic_share,
		c.owner_share as c_owner_share,
		c.method_type as c_method_type,
		c.is_active as c_is_active,
		c.with_holding_tax as c_with_holding_tax,
		c.created_at as c_created_at,
		c.updated_at as c_updated_at,
		c.deleted_at as c_deleted_at
	FROM tbl_user_clinic uc
	INNER JOIN tbl_clinic c ON uc.clinic_id = c.id
	WHERE uc.user_id = $1 AND uc.deleted_at IS NULL AND c.deleted_at IS NULL
	`
	var rows []user.UserClinicWithClinic
	if err := r.db.SelectContext(ctx, &rows, query, userID); err != nil {
		return nil, errors.New("failed to get user clinics")
	}
	if rows == nil {
		rows = []user.UserClinicWithClinic{}
	}
	out := make([]user.UserClinicWithClinic, 0, len(rows))
	for _, row := range rows {
		out = append(out, user.UserClinicWithClinic{
			UC_ID:            row.UC_ID,
			UC_UserID:        row.UC_UserID,
			UC_ClinicID:      row.UC_ClinicID,
			UC_Role:          row.UC_Role,
			UC_CreatedAt:     row.UC_CreatedAt,
			UC_UpdatedAt:     row.UC_UpdatedAt,
			C_ID:             row.C_ID,
			C_Name:           row.C_Name,
			C_ABNNumber:      row.C_ABNNumber,
			C_Address:        row.C_Address,
			C_City:           row.C_City,
			C_State:          row.C_State,
			C_Postcode:       row.C_Postcode,
			C_Phone:          row.C_Phone,
			C_Email:          row.C_Email,
			C_Website:        row.C_Website,
			C_LogoURL:        row.C_LogoURL,
			C_Description:    row.C_Description,
			C_ShareType:      row.C_ShareType,
			C_ClinicShare:    row.C_ClinicShare,
			C_OwnerShare:     row.C_OwnerShare,
			C_MethodType:     row.C_MethodType,
			C_IsActive:       row.C_IsActive,
			C_WithHoldingTax: row.C_WithHoldingTax,
			C_CreatedAt:      row.C_CreatedAt,
			C_UpdatedAt:      row.C_UpdatedAt,
			C_DeletedAt:      row.C_DeletedAt,
		})
	}
	if out == nil {
		out = []user.UserClinicWithClinic{}
	}
	return out, nil
}

func (r *userClinicRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE tbl_user_clinic SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
