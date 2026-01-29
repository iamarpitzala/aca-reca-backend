package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateExpenseType(ctx context.Context, db *sqlx.DB, expenseType *domain.ExpenseType) error {
	// Use NULL for deleted_at/deleted_by - uuid.Nil violates deleted_by foreign key
	query := `INSERT INTO tbl_expense_type (id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by)
		VALUES (:id, :clinic_id, :name, :description, :created_at, :created_by, NULL, NULL)`
	_, err := db.NamedExecContext(ctx, query, expenseType)
	if err != nil {
		return err
	}
	return nil
}

func CreateExpenseCategory(ctx context.Context, db *sqlx.DB, expenseCategory *domain.ExpenseCategory) error {
	// Use NULL for deleted_at/deleted_by - uuid.Nil violates deleted_by foreign key
	query := `INSERT INTO tbl_expense_category (id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by)
		VALUES (:id, :clinic_id, :name, :description, :created_at, :created_by, NULL, NULL)`
	_, err := db.NamedExecContext(ctx, query, expenseCategory)
	if err != nil {
		return err
	}
	return nil
}

func CreateExpenseCategoryType(ctx context.Context, db *sqlx.DB, expenseCategoryType *domain.ExpenseCategoryType) error {
	// Use NULL for deleted_at/deleted_by - uuid.Nil violates deleted_by foreign key
	query := `INSERT INTO tbl_expense_category_type (id, clinic_id, type_id, category_id, created_at, created_by, deleted_at, deleted_by)
		VALUES (:id, :clinic_id, :type_id, :category_id, :created_at, :created_by, NULL, NULL)`
	_, err := db.NamedExecContext(ctx, query, expenseCategoryType)
	if err != nil {
		return err
	}
	return nil
}

func CreateExpenseEntry(ctx context.Context, db *sqlx.DB, expenseEntry *domain.ExpenseEntry) error {
	// Use NULL for deleted_at/deleted_by - uuid.Nil violates deleted_by foreign key
	query := `INSERT INTO tbl_expense_entry (id, clinic_id, category_id, type_id, amount, gst_rate, is_gst_inclusive, expense_date, supplier_name, notes, created_at, created_by, deleted_at, deleted_by)
		VALUES (:id, :clinic_id, :category_id, :type_id, :amount, :gst_rate, :is_gst_inclusive, :expense_date, :supplier_name, :notes, :created_at, :created_by, NULL, NULL)`
	_, err := db.NamedExecContext(ctx, query, expenseEntry)
	if err != nil {
		return err
	}
	return nil
}

func GetExpenseTypesByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.ExpenseType, error) {
	query := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_type WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY name`
	var types []domain.ExpenseType
	err := db.SelectContext(ctx, &types, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to list expense types")
	}
	return types, nil
}

func GetExpenseTypeByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.ExpenseType, error) {
	query := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_type WHERE id = $1 AND deleted_at IS NULL`
	var expenseType domain.ExpenseType
	err := db.GetContext(ctx, &expenseType, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get expense type by id")
	}
	return &expenseType, nil
}

func GetExpenseCategoryByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.ExpenseCategory, error) {
	query := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category WHERE id = $1 AND deleted_at IS NULL`
	var expenseCategory domain.ExpenseCategory
	err := db.GetContext(ctx, &expenseCategory, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get expense category by id")
	}
	return &expenseCategory, nil
}

func GetExpenseCategoryTypeByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.ExpenseCategoryType, error) {
	query := `SELECT id, clinic_id, type_id, category_id, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category_type WHERE id = $1 AND deleted_at IS NULL`
	var expenseCategoryType domain.ExpenseCategoryType
	err := db.GetContext(ctx, &expenseCategoryType, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get expense category type by id")
	}
	return &expenseCategoryType, nil
}

func GetExpenseEntryByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.ExpenseEntry, error) {
	query := `SELECT id, clinic_id, category_id, type_id, amount, gst_rate, is_gst_inclusive, expense_date, supplier_name, notes, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_entry WHERE id = $1 AND deleted_at IS NULL`
	var expenseEntry domain.ExpenseEntry
	err := db.GetContext(ctx, &expenseEntry, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get expense entry by id")
	}
	return &expenseEntry, nil
}

func GetExpenseCategoriesByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.ExpenseCategory, error) {
	query := `SELECT id, clinic_id, name, description, created_at, created_by, deleted_at, deleted_by FROM tbl_expense_category WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY name`
	var categories []domain.ExpenseCategory
	err := db.SelectContext(ctx, &categories, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to list expense categories")
	}
	return categories, nil
}

func UpdateExpenseCategory(ctx context.Context, db *sqlx.DB, category *domain.ExpenseCategory) error {
	query := `UPDATE tbl_expense_category SET name = :name, description = :description WHERE id = :id AND deleted_at IS NULL`
	result, err := db.NamedExecContext(ctx, query, category)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("expense category not found")
	}
	return nil
}

func DeleteExpenseCategory(ctx context.Context, db *sqlx.DB, id uuid.UUID, deletedBy uuid.UUID) error {
	query := `UPDATE tbl_expense_category SET deleted_at = CURRENT_TIMESTAMP, deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := db.ExecContext(ctx, query, deletedBy, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("expense category not found")
	}
	return nil
}
