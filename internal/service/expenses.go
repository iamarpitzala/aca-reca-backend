package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type ExpensesService struct {
	db *sqlx.DB
}

func NewExpensesService(db *sqlx.DB) *ExpensesService {
	return &ExpensesService{
		db: db,
	}
}

func (es *ExpensesService) CreateExpenseType(ctx context.Context, expenseType *domain.ExpenseType) error {
	return repository.CreateExpenseType(ctx, es.db, expenseType)
}

func (es *ExpensesService) CreateExpenseCategory(ctx context.Context, expenseCategory *domain.ExpenseCategory) error {
	return repository.CreateExpenseCategory(ctx, es.db, expenseCategory)
}

func (es *ExpensesService) CreateExpenseCategoryType(ctx context.Context, expenseCategoryType *domain.ExpenseCategoryType) error {
	return repository.CreateExpenseCategoryType(ctx, es.db, expenseCategoryType)
}

func (es *ExpensesService) CreateExpenseEntry(ctx context.Context, expenseEntry *domain.ExpenseEntry) error {
	return repository.CreateExpenseEntry(ctx, es.db, expenseEntry)
}

func (es *ExpensesService) GetExpenseTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseType, error) {
	return repository.GetExpenseTypeByID(ctx, es.db, id)
}

func (es *ExpensesService) GetExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	return repository.GetExpenseCategoryByID(ctx, es.db, id)
}

func (es *ExpensesService) GetExpenseCategoryTypeByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategoryType, error) {
	return repository.GetExpenseCategoryTypeByID(ctx, es.db, id)
}

func (es *ExpensesService) GetExpenseEntryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseEntry, error) {
	return repository.GetExpenseEntryByID(ctx, es.db, id)
}
