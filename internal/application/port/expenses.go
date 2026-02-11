package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ExpenseRepository interface {
	CreateExpenseType(ctx context.Context, t *domain.ExpenseType) error
	CreateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error
	CreateExpenseCategoryType(ctx context.Context, ct *domain.ExpenseCategoryType) error
	CreateExpenseEntry(ctx context.Context, e *domain.ExpenseEntry) error
	GetExpenseTypesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseType, error)
	GetExpenseTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseType, error)
	GetExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error)
	GetExpenseCategoryTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategoryType, error)
	GetExpenseEntryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseEntry, error)
	GetExpenseEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseEntry, error)
	GetExpenseCategoriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ExpenseCategory, error)
	UpdateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error
	DeleteExpenseCategory(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}
