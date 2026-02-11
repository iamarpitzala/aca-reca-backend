package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type expenseRepo struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) port.ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) CreateExpenseType(ctx context.Context, t *domain.ExpenseType) error {
	q := `INSERT INTO tbl_expense_type (id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by) VALUES (:id, :clinic_id, :name, :description, :created_at, :created_by, NULL, NULL)`
	_, err := r.db.NamedExecContext(ctx, q, t)
	return err
}

func (r *expenseRepo) CreateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	q := `INSERT INTO tbl_expense_category (id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by) VALUES (:id, :clinic_id, :name, :description, :created_at, :created_by, NULL, NULL)`
	_, err := r.db.NamedExecContext(ctx, q, c)
	return err
}

func (r *expenseRepo) CreateExpenseCategoryType(ctx context.Context, ct *domain.ExpenseCategoryType) error {
	q := `INSERT INTO tbl_expense_category_type (id, clinic_id, type_id, category_id, created_at, created_by, deleted_at, deleted_by) VALUES (:id, :clinic_id, :type_id, :category_id, :created_at, :created_by, NULL, NULL)`
	_, err := r.db.NamedExecContext(ctx, q, ct)
	return err
}

func (r *expenseRepo) CreateExpenseEntry(ctx context.Context, e *domain.ExpenseEntry) error {
	q := `INSERT INTO tbl_expense_entry (id, clinic_id, category_id, type_id, amount, gst_rate, is_gst_inclusive, expense_date, supplier_name, notes, created_at, created_by, deleted_at, deleted_by) VALUES (:id, :clinic_id, :category_id, :type_id, :amount, :gst_rate, :is_gst_inclusive, :expense_date, :supplier_name, :notes, :created_at, :created_by, NULL, NULL)`
	_, err := r.db.NamedExecContext(ctx, q, e)
	return err
}

func (r *expenseRepo) GetExpenseTypesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseType, error) {
	q := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_type WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY name`
	var out []domain.ExpenseType
	if err := r.db.SelectContext(ctx, &out, q, clinicID); err != nil {
		return nil, errors.New("failed to list expense types")
	}
	return out, nil
}

func (r *expenseRepo) GetExpenseTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseType, error) {
	q := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_type WHERE id = $1 AND deleted_at IS NULL`
	var t domain.ExpenseType
	if err := r.db.GetContext(ctx, &t, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get expense type by id")
	}
	return &t, nil
}

func (r *expenseRepo) GetExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	q := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category WHERE id = $1 AND deleted_at IS NULL`
	var c domain.ExpenseCategory
	if err := r.db.GetContext(ctx, &c, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get expense category by id")
	}
	return &c, nil
}

func (r *expenseRepo) GetExpenseCategoryTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategoryType, error) {
	q := `SELECT id, clinic_id, type_id, category_id, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category_type WHERE id = $1 AND deleted_at IS NULL`
	var ct domain.ExpenseCategoryType
	if err := r.db.GetContext(ctx, &ct, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get expense category type by id")
	}
	return &ct, nil
}

func (r *expenseRepo) GetExpenseEntryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseEntry, error) {
	q := `SELECT id, clinic_id, category_id, type_id, amount, gst_rate, is_gst_inclusive, expense_date, supplier_name, notes, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_entry WHERE id = $1 AND deleted_at IS NULL`
	var e domain.ExpenseEntry
	if err := r.db.GetContext(ctx, &e, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get expense entry by id")
	}
	return &e, nil
}

func (r *expenseRepo) GetExpenseEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseEntry, error) {
	q := `SELECT id, clinic_id, category_id, type_id, amount, gst_rate, is_gst_inclusive, expense_date, supplier_name, notes, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_entry WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY expense_date DESC`
	var out []domain.ExpenseEntry
	if err := r.db.SelectContext(ctx, &out, q, clinicID); err != nil {
		return nil, errors.New("failed to list expense entries for clinic")
	}
	return out, nil
}

func (r *expenseRepo) GetExpenseCategoriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseCategory, error) {
	q := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY name`
	var out []domain.ExpenseCategory
	if err := r.db.SelectContext(ctx, &out, q, clinicID); err != nil {
		return nil, errors.New("failed to list expense categories")
	}
	return out, nil
}

func (r *expenseRepo) UpdateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	res, err := r.db.NamedExecContext(ctx, `UPDATE tbl_expense_category SET name = :name, description = :description WHERE id = :id AND deleted_at IS NULL`, c)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("expense category not found")
	}
	return nil
}

func (r *expenseRepo) DeleteExpenseCategory(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_expense_category SET deleted_at = CURRENT_TIMESTAMP, deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("expense category not found")
	}
	return nil
}
