package usecase

import (
	"context"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
)

type COAService struct {
	repo port.ChartOfAccountsRepository
}

func NewCOAService(repo port.ChartOfAccountsRepository) *COAService {
	return &COAService{repo: repo}
}

func (s *COAService) CreateCOA(ctx context.Context, req *coa.COARequest, ownerUserID string) error {
	coa := req.ToRepo()
	coa.OwnerUserID = ownerUserID
	existing, err := s.repo.GetByCodeAndOwner(ctx, coa.Code, ownerUserID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("chart of accounts with this code already exists")
	}
	accountType, err := s.repo.GetAccountTypeByID(ctx, coa.AccountTypeID)
	if err != nil {
		return err
	}
	if accountType == nil {
		return errors.New("account type not found")
	}
	accountTax, err := s.repo.GetAccountTaxByID(ctx, coa.AccountTaxID)
	if err != nil {
		return err
	}
	if accountTax == nil {
		return errors.New("account tax not found")
	}

	return s.repo.Create(ctx, coa)
}

func (s *COAService) GetCOAByID(ctx context.Context, id string) (*coa.COAResponse, error) {
	coa, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, errors.New("coa not found")
	}
	return coa.ToResponse(), nil
}

func (s *COAService) GetCOAByCode(ctx context.Context, code string) (*coa.COAResponse, error) {
	coa, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, errors.New("coa not found")
	}
	return coa.ToResponse(), nil
}

func (s *COAService) GetCOAByAccountTypeID(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]coa.COAResponse, error) {
	coas, err := s.repo.GetByAccountTypeIDSorted(ctx, accountTypeID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	out := make([]coa.COAResponse, 0, len(coas))
	for _, coa := range coas {
		out = append(out, *coa.ToResponse())
	}
	return out, nil
}

func (s *COAService) GetCOAsByAccountType(ctx context.Context, sortBy, sortOrder string) ([]coa.COAResponse, error) {
	coas, err := s.repo.GetByAccountTypeSorted(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	out := make([]coa.COAResponse, 0, len(coas))
	for _, coa := range coas {
		out = append(out, *coa.ToResponse())
	}
	return out, nil
}

func (s *COAService) GetCOAByAccountTaxID(ctx context.Context, accountTaxID int) ([]coa.COAResponse, error) {
	coas, err := s.repo.GetByAccountTaxID(ctx, accountTaxID)
	if err != nil {
		return nil, err
	}
	out := make([]coa.COAResponse, 0, len(coas))
	for _, coa := range coas {
		out = append(out, *coa.ToResponse())
	}
	return out, nil
}

func (s *COAService) GetAllCOAs(ctx context.Context) ([]coa.COAResponse, error) {
	coas, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]coa.COAResponse, 0, len(coas))
	for _, coa := range coas {
		out = append(out, *coa.ToResponse())
	}
	return out, nil
}

func (s *COAService) UpdateCOA(ctx context.Context, coa *coa.COA) error {
	existing, err := s.repo.GetByID(ctx, coa.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("coa not found")
	}
	existing.AccountTypeID = coa.AccountTypeID
	existing.AccountTaxID = coa.AccountTaxID
	existing.Code = coa.Code
	existing.Name = coa.Name
	existing.Description = coa.Description
	return s.repo.Update(ctx, existing)
}

func (s *COAService) DeleteCOA(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.Delete(ctx, ids)
}

func (s *COAService) BulkUpdateCOATax(ctx context.Context, ids []string, accountTaxID int) error {
	if len(ids) == 0 {
		return nil
	}
	tax, err := s.repo.GetAccountTaxByID(ctx, accountTaxID)
	if err != nil || tax == nil {
		return errors.New("account tax not found")
	}
	return s.repo.BulkUpdateAccountTax(ctx, ids, accountTaxID)
}

func (s *COAService) GetCOAType(ctx context.Context) ([]coa.AccountTypeCOA, error) {
	return s.repo.GetAllAccountTypes(ctx)
}

func (s *COAService) GetAccountTax(ctx context.Context) ([]coa.AccountTaxCOA, error) {
	return s.repo.GetAllAccountTax(ctx)
}

// CreateDefaultAccountsForUser creates the default chart of accounts for a user (same set as seed).
// Call this when a new user is created so they get the standard accounts automatically.
func (s *COAService) CreateDefaultAccountsForUser(ctx context.Context, userID string) error {
	return s.repo.CreateDefaultAccountsForUser(ctx, userID)
}

func (s *COAService) CheckIfAccountTaxIsTaxable(ctx context.Context, accountTypeID int) (bool, error) {
	return s.repo.CheckIfAccountTaxIsTaxable(ctx, accountTypeID)
}
