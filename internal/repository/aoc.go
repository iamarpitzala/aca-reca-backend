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

func CreateAOC(ctx context.Context, db *sqlx.DB, aoc *domain.AOC) error {
	query := `INSERT INTO tbl_aoc (id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at)
		VALUES (:id, :account_type_id, :account_tax_id, :code, :name, :description, :created_at, :updated_at, :deleted_at)`
	_, err := db.NamedExecContext(ctx, query, aoc)
	if err != nil {
		return err
	}
	return nil
}

func GetAOCByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_aoc WHERE id = $1 AND deleted_at IS NULL`
	var aoc domain.AOC
	err := db.GetContext(ctx, &aoc, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get aoc by id")
	}
	return &aoc, nil
}

func UpdateAOC(ctx context.Context, db *sqlx.DB, aoc *domain.AOC) error {
	query := `UPDATE tbl_aoc SET account_type_id = :account_type_id, account_tax_id = :account_tax_id, code = :code, name = :name, description = :description, updated_at = :updated_at WHERE id = :id`
	_, err := db.NamedExecContext(ctx, query, aoc)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAOC(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_aoc SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllAOCs(ctx context.Context, db *sqlx.DB) ([]domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_aoc WHERE deleted_at IS NULL`
	var aocs []domain.AOC
	err := db.SelectContext(ctx, &aocs, query)
	if err != nil {
		return nil, errors.New("failed to get all aocs")
	}
	return aocs, nil
}

func GetAOCByCode(ctx context.Context, db *sqlx.DB, code string) (*domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_aoc WHERE code = $1 AND deleted_at IS NULL`
	var aoc domain.AOC
	err := db.GetContext(ctx, &aoc, query, code)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get aoc by code")
	}
	return &aoc, nil
}

func GetAOCByAccountTypeID(ctx context.Context, db *sqlx.DB, accountTypeID int) ([]domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_aoc WHERE account_type_id = $1 AND deleted_at IS NULL`
	var aocs []domain.AOC
	err := db.SelectContext(ctx, &aocs, query, accountTypeID)
	if err != nil {
		return nil, errors.New("failed to get aocs by account type id")
	}
	return aocs, nil
}

func GetAOCByAccountTaxID(ctx context.Context, db *sqlx.DB, accountTaxID int) ([]domain.AOC, error) {
	query := `SELECT id, account_type_id, account_tax_id, code, name, description, created_at, updated_at, deleted_at FROM tbl_aoc WHERE account_tax_id = $1 AND deleted_at IS NULL`
	var aocs []domain.AOC
	err := db.SelectContext(ctx, &aocs, query, accountTaxID)
	if err != nil {
		return nil, errors.New("failed to get aocs by account tax id")
	}
	return aocs, nil
}

func GetAllAOCType(ctx context.Context, db *sqlx.DB) ([]domain.AccountType, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_account_type WHERE deleted_at IS NULL`
	var accountTypes []domain.AccountType
	err := db.SelectContext(ctx, &accountTypes, query)
	if err != nil {
		return nil, errors.New("failed to get all account types")
	}
	return accountTypes, nil
}
