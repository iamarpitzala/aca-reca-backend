package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func GetAOCByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE id = $1 AND deleted_at IS NULL`
	var aoc domain.AOC
	err := db.GetContext(ctx, &aoc, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &aoc, nil
}

func GetAccountTaxByID(ctx context.Context, db *sqlx.DB, id int) (*domain.AccountTax, error) {
	query := `SELECT id, name, rate, description, created_at, updated_at FROM tbl_account_tax WHERE id = $1 AND deleted_at IS NULL`
	var at domain.AccountTax
	err := db.GetContext(ctx, &at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &at, nil
}

// TaxNameToCategory maps tbl_account_tax.name to frontend tax category
func TaxNameToCategory(name string) string {
	switch name {
	case "GST on Income":
		return domain.TaxCategoryGSTOnIncome
	case "GST Free Income":
		return domain.TaxCategoryGSTFreeIncome
	case "GST on Expenses":
		return domain.TaxCategoryGSTOnExpenses
	case "GST Free Expenses":
		return domain.TaxCategoryGSTFreeExpenses
	case "BAS Excluded":
		return domain.TaxCategoryBASExcluded
	default:
		return domain.TaxCategoryBASExcluded
	}
}

// ClinicCOAAssigned returns true if the clinic has the given COA assigned
func ClinicCOAAssigned(ctx context.Context, db *sqlx.DB, clinicID, coaID uuid.UUID) (bool, error) {
	var n int
	query := `SELECT 1 FROM tbl_clinic_coa WHERE clinic_id = $1 AND coa_id = $2 AND deleted_at IS NULL LIMIT 1`
	err := db.GetContext(ctx, &n, query, clinicID, coaID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListClinicCOAs returns AOC records assigned to the clinic (for mapping dropdown)
func ListClinicCOAs(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.AOC, error) {
	query := `SELECT a.id, a.account_type_id, a.account_tax_id, a.code, a.name, a.description, a.created_at, a.updated_at, a.deleted_at
		FROM tbl_account a
		INNER JOIN tbl_clinic_coa cc ON cc.coa_id = a.id AND cc.clinic_id = $1 AND cc.deleted_at IS NULL
		WHERE a.deleted_at IS NULL
		ORDER BY a.code`
	var list []domain.AOC
	err := db.SelectContext(ctx, &list, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to list clinic COAs")
	}
	if list == nil {
		list = []domain.AOC{}
	}
	return list, nil
}
