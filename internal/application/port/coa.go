package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
)

type ChartOfAccountsRepository interface {
	Create(ctx context.Context, coa *coa.COA) error
	GetByID(ctx context.Context, id string) (*coa.COA, error)
	Update(ctx context.Context, coa *coa.COA) error
	Delete(ctx context.Context, ids []string) error
	BulkUpdateAccountTax(ctx context.Context, ids []string, accountTaxID int) error
	List(ctx context.Context) ([]coa.COA, error)
	GetByCode(ctx context.Context, code string) (*coa.COA, error)
	GetByCodeAndOwner(ctx context.Context, code, ownerUserID string) (*coa.COA, error)
	GetByAccountTypeID(ctx context.Context, accountTypeID int) ([]coa.COA, error)
	GetByAccountTypeIDSorted(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]coa.COA, error)
	GetByAccountTypeSorted(ctx context.Context, sortBy, sortOrder string) ([]coa.COA, error)
	CheckIfAccountTaxIsTaxable(ctx context.Context, accountTypeID int) (bool, error)
	GetByAccountTaxID(ctx context.Context, accountTaxID int) ([]coa.COA, error)
	GetAllAccountTypes(ctx context.Context) ([]coa.AccountTypeCOA, error)
	GetAccountTypeByID(ctx context.Context, id int) (*coa.AccountTypeCOA, error)
	GetAllAccountTax(ctx context.Context) ([]coa.AccountTaxCOA, error)
	GetAccountTaxByID(ctx context.Context, id int) (*coa.AccountTaxCOA, error)
	CreateDefaultAccountsForUser(ctx context.Context, userID string) error
}
