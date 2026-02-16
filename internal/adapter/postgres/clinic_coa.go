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

type clinicCOAWithDetailsRow struct {
	ID              uuid.UUID  `db:"id"`
	ClinicID        uuid.UUID  `db:"clinic_id"`
	COAID           uuid.UUID  `db:"coa_id"`
	Code            string     `db:"code"`
	Name            string     `db:"name"`
	AccountTypeID   int        `db:"account_type_id"`
	AccountTaxID    int        `db:"account_tax_id"`
	AccountTypeName *string    `db:"account_type_name"`
	AccountTaxName  *string    `db:"account_tax_name"`
	Description     *string    `db:"description"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

func (r *clinicCOARepo) ListByClinicIDWithDetails(ctx context.Context, clinicID uuid.UUID) ([]domain.ClinicCOAWithDetails, error) {
	query := `SELECT cc.id, cc.clinic_id, cc.coa_id, cc.created_at, cc.updated_at,
		a.code, a.name, a.account_type_id, a.account_tax_id, a.description,
		at.name AS account_type_name, ax.name AS account_tax_name
		FROM tbl_clinic_coa cc
		INNER JOIN tbl_account a ON a.id = cc.coa_id AND a.deleted_at IS NULL
		LEFT JOIN tbl_account_type at ON at.id = a.account_type_id
		LEFT JOIN tbl_account_tax ax ON ax.id = a.account_tax_id
		WHERE cc.clinic_id = $1 AND cc.deleted_at IS NULL
		ORDER BY a.code ASC`
	var rows []clinicCOAWithDetailsRow
	if err := r.db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, err
	}
	if rows == nil {
		return []domain.ClinicCOAWithDetails{}, nil
	}
	out := make([]domain.ClinicCOAWithDetails, len(rows))
	for i := range rows {
		out[i] = domain.ClinicCOAWithDetails{
			ID:              rows[i].ID,
			ClinicID:        rows[i].ClinicID,
			COAID:           rows[i].COAID,
			Code:            rows[i].Code,
			Name:            rows[i].Name,
			AccountTypeID:   rows[i].AccountTypeID,
			AccountTaxID:    rows[i].AccountTaxID,
			AccountTypeName: safeStr(rows[i].AccountTypeName),
			AccountTaxName:  safeStr(rows[i].AccountTaxName),
			Description:     rows[i].Description,
			CreatedAt:       rows[i].CreatedAt,
			UpdatedAt:       rows[i].UpdatedAt,
		}
	}
	return out, nil
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
