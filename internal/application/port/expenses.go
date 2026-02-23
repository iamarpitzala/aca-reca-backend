package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ExpenseRepository interface {
	CreateExpenseType(ctx context.Context, t *domain.ExpenseType) error
	CreateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error
	CreateExpenseCategoryType(ctx context.Context, ct *domain.ExpenseCategoryType) error
	CreateExpenseEntry(ctx context.Context, e *domain.ExpenseEntry) error
	GetExpenseTypesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseType, error)
	GetExpenseTypeByID(ctx context.Context, id string) (*domain.ExpenseType, error)
	GetExpenseCategoryByID(ctx context.Context, id string) (*domain.ExpenseCategory, error)
	GetExpenseCategoryTypeByID(ctx context.Context, id string) (*domain.ExpenseCategoryType, error)
	GetExpenseEntryByID(ctx context.Context, id string) (*domain.ExpenseEntry, error)
	GetExpenseEntriesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseEntry, error)
	GetExpenseCategoriesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseCategory, error)
	UpdateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error
	DeleteExpenseCategory(ctx context.Context, id string, deletedBy string) error
}
