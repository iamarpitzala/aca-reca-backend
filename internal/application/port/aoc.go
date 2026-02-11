package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type AOCRepository interface {
	Create(ctx context.Context, aoc *domain.AOC) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AOC, error)
	Update(ctx context.Context, aoc *domain.AOC) error
	Delete(ctx context.Context, ids []uuid.UUID) error
	BulkUpdateAccountTax(ctx context.Context, ids []uuid.UUID, accountTaxID int) error
	List(ctx context.Context) ([]domain.AOC, error)
	GetByCode(ctx context.Context, code string) (*domain.AOC, error)
	GetByAccountTypeID(ctx context.Context, accountTypeID int) ([]domain.AOC, error)
	GetByAccountTypeIDSorted(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]domain.AOC, error)
	GetByAccountTypeSorted(ctx context.Context, sortBy, sortOrder string) ([]domain.AOC, error)
	GetByAccountTaxID(ctx context.Context, accountTaxID int) ([]domain.AOC, error)
	GetAllAccountTypes(ctx context.Context) ([]domain.AccountType, error)
	GetAccountTypeByID(ctx context.Context, id int) (*domain.AccountType, error)
	GetAllAccountTax(ctx context.Context) ([]domain.AccountTax, error)
	GetAccountTaxByID(ctx context.Context, id int) (*domain.AccountTax, error)
}
