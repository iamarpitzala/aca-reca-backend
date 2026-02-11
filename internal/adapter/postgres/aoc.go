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

type aocRepo struct {
	db *sqlx.DB
}

func NewAOCRepository(db *sqlx.DB) port.AOCRepository {
	return &aocRepo{db: db}
}

func (r *aocRepo) Create(ctx context.Context, aoc *domain.AOC) error {
	query := `INSERT INTO tbl_account (id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at) VALUES (:id, :account_type_id, :account_tax_id, :code, :name, :description, :created_at, :updated_at, :deleted_at)`
	_, err := r.db.NamedExecContext(ctx, query, aoc)
	return err
}

func (r *aocRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE id = $1 AND deleted_at IS NULL`
	var aoc domain.AOC
	err := r.db.GetContext(ctx, &aoc, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get aoc by id")
	}
	return &aoc, nil
}

func (r *aocRepo) Update(ctx context.Context, aoc *domain.AOC) error {
	query := `UPDATE tbl_account SET account_type_id = :account_type_id, account_tax_id = :account_tax_id, code = :code, name = :name, description = :description, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, aoc)
	return err
}

func (r *aocRepo) Delete(ctx context.Context, ids []uuid.UUID) error {
	now := time.Now()
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx, `UPDATE tbl_account SET deleted_at = $1 WHERE id = $2`, now, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *aocRepo) BulkUpdateAccountTax(ctx context.Context, ids []uuid.UUID, accountTaxID int) error {
	now := time.Now()
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx, `UPDATE tbl_account SET account_tax_id = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`, accountTaxID, now, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *aocRepo) List(ctx context.Context) ([]domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE deleted_at IS NULL ORDER BY code`
	var aocs []domain.AOC
	err := r.db.SelectContext(ctx, &aocs, query)
	if err != nil {
		return nil, err
	}
	if aocs == nil {
		aocs = []domain.AOC{}
	}
	return aocs, nil
}

func (r *aocRepo) GetByCode(ctx context.Context, code string) (*domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE code = $1 AND deleted_at IS NULL`
	var aoc domain.AOC
	err := r.db.GetContext(ctx, &aoc, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get aoc by code")
	}
	return &aoc, nil
}

func (r *aocRepo) GetByAccountTypeID(ctx context.Context, accountTypeID int) ([]domain.AOC, error) {
	return r.GetByAccountTypeIDSorted(ctx, accountTypeID, "code", "asc")
}

func (r *aocRepo) GetByAccountTypeIDSorted(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]domain.AOC, error) {
	col := "code"
	if sortBy == "name" {
		col = "name"
	}
	dir := "ASC"
	if sortOrder == "desc" {
		dir = "DESC"
	}
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE account_type_id = $1 AND deleted_at IS NULL ORDER BY ` + col + " " + dir
	var aocs []domain.AOC
	err := r.db.SelectContext(ctx, &aocs, query, accountTypeID)
	if err != nil {
		return nil, err
	}
	if aocs == nil {
		aocs = []domain.AOC{}
	}
	return aocs, nil
}

func (r *aocRepo) GetByAccountTypeSorted(ctx context.Context, sortBy, sortOrder string) ([]domain.AOC, error) {
	col := "code"
	if sortBy == "name" {
		col = "name"
	}
	dir := "ASC"
	if sortOrder == "desc" {
		dir = "DESC"
	}
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE deleted_at IS NULL ORDER BY ` + col + " " + dir
	var aocs []domain.AOC
	err := r.db.SelectContext(ctx, &aocs, query)
	if err != nil {
		return nil, err
	}
	if aocs == nil {
		aocs = []domain.AOC{}
	}
	return aocs, nil
}

func (r *aocRepo) GetByAccountTaxID(ctx context.Context, accountTaxID int) ([]domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_account WHERE account_tax_id = $1 AND deleted_at IS NULL`
	var aocs []domain.AOC
	err := r.db.SelectContext(ctx, &aocs, query, accountTaxID)
	if err != nil {
		return nil, errors.New("failed to get aocs by account tax id")
	}
	return aocs, nil
}

func (r *aocRepo) GetAllAccountTypes(ctx context.Context) ([]domain.AccountType, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_account_type WHERE deleted_at IS NULL`
	var out []domain.AccountType
	err := r.db.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, errors.New("failed to get all account types")
	}
	return out, nil
}

func (r *aocRepo) GetAccountTypeByID(ctx context.Context, id int) (*domain.AccountType, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_account_type WHERE id = $1 AND deleted_at IS NULL`
	var at domain.AccountType
	err := r.db.GetContext(ctx, &at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get account type by id")
	}
	return &at, nil
}

func (r *aocRepo) GetAllAccountTax(ctx context.Context) ([]domain.AccountTax, error) {
	query := `SELECT id, name, rate, description, created_at, updated_at FROM tbl_account_tax WHERE deleted_at IS NULL`
	var out []domain.AccountTax
	err := r.db.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, errors.New("failed to get all account taxes")
	}
	return out, nil
}

func (r *aocRepo) GetAccountTaxByID(ctx context.Context, id int) (*domain.AccountTax, error) {
	query := `SELECT id, name, rate, description, created_at, updated_at FROM tbl_account_tax WHERE id = $1 AND deleted_at IS NULL`
	var at domain.AccountTax
	err := r.db.GetContext(ctx, &at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get account tax by id")
	}
	return &at, nil
}
