package usecase

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ExpensesService struct {
	repo port.ExpenseRepository
}

func NewExpensesService(repo port.ExpenseRepository) *ExpensesService {
	return &ExpensesService{repo: repo}
}

func (s *ExpensesService) CreateExpenseType(ctx context.Context, t *domain.ExpenseType) error {
	return s.repo.CreateExpenseType(ctx, t)
}

func (s *ExpensesService) CreateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	return s.repo.CreateExpenseCategory(ctx, c)
}

func (s *ExpensesService) CreateExpenseCategoryType(ctx context.Context, ct *domain.ExpenseCategoryType) error {
	return s.repo.CreateExpenseCategoryType(ctx, ct)
}

func (s *ExpensesService) CreateExpenseEntry(ctx context.Context, e *domain.ExpenseEntry) error {
	return s.repo.CreateExpenseEntry(ctx, e)
}

func (s *ExpensesService) GetExpenseTypesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseType, error) {
	return s.repo.GetExpenseTypesByClinicID(ctx, clinicID)
}

func (s *ExpensesService) GetExpenseTypeByID(ctx context.Context, id string) (*domain.ExpenseType, error) {
	return s.repo.GetExpenseTypeByID(ctx, id)
}

func (s *ExpensesService) GetExpenseCategoryByID(ctx context.Context, id string) (*domain.ExpenseCategory, error) {
	return s.repo.GetExpenseCategoryByID(ctx, id)
}

func (s *ExpensesService) GetExpenseCategoryTypeByID(ctx context.Context, id string) (*domain.ExpenseCategoryType, error) {
	return s.repo.GetExpenseCategoryTypeByID(ctx, id)
}

func (s *ExpensesService) GetExpenseEntryByID(ctx context.Context, id string) (*domain.ExpenseEntry, error) {
	return s.repo.GetExpenseEntryByID(ctx, id)
}

func (s *ExpensesService) GetExpenseEntriesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseEntry, error) {
	return s.repo.GetExpenseEntriesByClinicID(ctx, clinicID)
}

func (s *ExpensesService) GetExpenseCategoriesByClinicID(ctx context.Context, clinicID string) ([]domain.ExpenseCategory, error) {
	return s.repo.GetExpenseCategoriesByClinicID(ctx, clinicID)
}

func (s *ExpensesService) UpdateExpenseCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	return s.repo.UpdateExpenseCategory(ctx, c)
}

func (s *ExpensesService) DeleteExpenseCategory(ctx context.Context, id string, deletedBy string) error {
	return s.repo.DeleteExpenseCategory(ctx, id, deletedBy)
}
