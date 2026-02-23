package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
	"github.com/jmoiron/sqlx"
)

type coaRepo struct {
	db *sqlx.DB
}

func NewChartOfAccountsRepository(db *sqlx.DB) port.ChartOfAccountsRepository {
	return &coaRepo{db: db}
}

func (r *coaRepo) Create(ctx context.Context, coa *coa.COA) error {
	query := `INSERT INTO tbl_account (id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at) VALUES (:id, :account_type_id, :account_tax_id, :code, :name, :description, :created_at, :updated_at, :deleted_at)`
	_, err := r.db.NamedExecContext(ctx, query, coa)
	return err
}

func (r *coaRepo) GetByID(ctx context.Context, id string) (*coa.COA, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE id = $1 AND deleted_at IS NULL`
	var coa coa.COA
	err := r.db.GetContext(ctx, &coa, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get coa by id")
	}
	return &coa, nil
}

func (r *coaRepo) Update(ctx context.Context, coa *coa.COA) error {
	query := `UPDATE tbl_account SET account_type_id = :account_type_id, account_tax_id = :account_tax_id, code = :code, name = :name, description = :description, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, coa)
	return err
}

func (r *coaRepo) Delete(ctx context.Context, ids []string) error {
	now := time.Now()
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx, `UPDATE tbl_account SET deleted_at = $1 WHERE id = $2`, now, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *coaRepo) BulkUpdateAccountTax(ctx context.Context, ids []string, accountTaxID int) error {
	now := time.Now()
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx, `UPDATE tbl_account SET account_tax_id = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`, accountTaxID, now, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *coaRepo) List(ctx context.Context) ([]coa.COA, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE deleted_at IS NULL ORDER BY code`
	var coas []coa.COA
	err := r.db.SelectContext(ctx, &coas, query)
	if err != nil {
		return nil, err
	}
	if coas == nil {
		coas = []coa.COA{}
	}
	return coas, nil
}

func (r *coaRepo) GetByCode(ctx context.Context, code string) (*coa.COA, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE code = $1 AND deleted_at IS NULL`
	var coa coa.COA
	err := r.db.GetContext(ctx, &coa, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get coa by code")
	}
	return &coa, nil
}

func (r *coaRepo) GetByAccountTypeID(ctx context.Context, accountTypeID int) ([]coa.COA, error) {
	return r.GetByAccountTypeIDSorted(ctx, accountTypeID, "code", "asc")
}

func (r *coaRepo) GetByAccountTypeIDSorted(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]coa.COA, error) {
	col := "code"
	if sortBy == "name" {
		col = "name"
	}
	dir := "ASC"
	if sortOrder == "desc" {
		dir = "DESC"
	}
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE account_type_id = $1 AND deleted_at IS NULL ORDER BY ` + col + " " + dir
	var coas []coa.COA
	err := r.db.SelectContext(ctx, &coas, query, accountTypeID)
	if err != nil {
		return nil, err
	}
	if coas == nil {
		coas = []coa.COA{}
	}
	return coas, nil
}

func (r *coaRepo) GetByAccountTypeSorted(ctx context.Context, sortBy, sortOrder string) ([]coa.COA, error) {
	col := "code"
	if sortBy == "name" {
		col = "name"
	}
	dir := "ASC"
	if sortOrder == "desc" {
		dir = "DESC"
	}
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE deleted_at IS NULL ORDER BY ` + col + " " + dir
	var coas []coa.COA
	err := r.db.SelectContext(ctx, &coas, query)
	if err != nil {
		return nil, err
	}
	if coas == nil {
		coas = []coa.COA{}
	}
	return coas, nil
}

func (r *coaRepo) GetByAccountTaxID(ctx context.Context, accountTaxID int) ([]coa.COA, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE account_tax_id = $1 AND deleted_at IS NULL`
	var coas []coa.COA
	err := r.db.SelectContext(ctx, &coas, query, accountTaxID)
	if err != nil {
		return nil, errors.New("failed to get coas by account tax id")
	}
	return coas, nil
}

func (r *coaRepo) GetAllAccountTypes(ctx context.Context) ([]coa.AccountTypeCOA, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_account_type WHERE deleted_at IS NULL`
	var out []coa.AccountTypeCOA
	err := r.db.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, errors.New("failed to get all account types")
	}
	return out, nil
}

func (r *coaRepo) GetAccountTypeByID(ctx context.Context, id int) (*coa.AccountTypeCOA, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_account_type WHERE id = $1 AND deleted_at IS NULL`
	var at coa.AccountTypeCOA
	err := r.db.GetContext(ctx, &at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get account type by id")
	}
	return &at, nil
}

func (r *coaRepo) GetAllAccountTax(ctx context.Context) ([]coa.AccountTaxCOA, error) {
	query := `SELECT id, name, rate, description, created_at, updated_at FROM tbl_account_tax WHERE deleted_at IS NULL`
	var out []coa.AccountTaxCOA
	err := r.db.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, errors.New("failed to get all account taxes")
	}
	return out, nil
}

func (r *coaRepo) GetAccountTaxByID(ctx context.Context, id int) (*coa.AccountTaxCOA, error) {
	query := `SELECT id, name, rate, description, created_at, updated_at FROM tbl_account_tax WHERE id = $1 AND deleted_at IS NULL`
	var at coa.AccountTaxCOA
	err := r.db.GetContext(ctx, &at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get account tax by id")
	}
	return &at, nil
}
