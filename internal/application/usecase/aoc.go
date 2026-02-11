package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type AOCService struct {
	repo port.AOCRepository
}

func NewAOCService(repo port.AOCRepository) *AOCService {
	return &AOCService{repo: repo}
}

func (s *AOCService) CreateAOC(ctx context.Context, req *domain.AOCRequest) error {
	aoc := req.ToRepo()
	existing, err := s.repo.GetByCode(ctx, aoc.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("aoc with this code already exists")
	}
	accountType, err := s.repo.GetAccountTypeByID(ctx, aoc.AccountTypeID)
	if err != nil {
		return err
	}
	if accountType == nil {
		return errors.New("account type not found")
	}
	accountTax, err := s.repo.GetAccountTaxByID(ctx, aoc.AccountTaxID)
	if err != nil {
		return err
	}
	if accountTax == nil {
		return errors.New("account tax not found")
	}
	return s.repo.Create(ctx, aoc)
}

func (s *AOCService) GetAOCByID(ctx context.Context, id uuid.UUID) (*domain.AOCResponse, error) {
	aoc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if aoc == nil {
		return nil, errors.New("aoc not found")
	}
	return aoc.ToResponse(), nil
}

func (s *AOCService) GetAOCByCode(ctx context.Context, code string) (*domain.AOCResponse, error) {
	aoc, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if aoc == nil {
		return nil, errors.New("aoc not found")
	}
	return aoc.ToResponse(), nil
}

func (s *AOCService) GetAOCByAccountTypeID(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]domain.AOCResponse, error) {
	aocs, err := s.repo.GetByAccountTypeIDSorted(ctx, accountTypeID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		out = append(out, *aoc.ToResponse())
	}
	return out, nil
}

func (s *AOCService) GetAOCsByAccountType(ctx context.Context, sortBy, sortOrder string) ([]domain.AOCResponse, error) {
	aocs, err := s.repo.GetByAccountTypeSorted(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		out = append(out, *aoc.ToResponse())
	}
	return out, nil
}

func (s *AOCService) GetAOCByAccountTaxID(ctx context.Context, accountTaxID int) ([]domain.AOCResponse, error) {
	aocs, err := s.repo.GetByAccountTaxID(ctx, accountTaxID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		out = append(out, *aoc.ToResponse())
	}
	return out, nil
}

func (s *AOCService) GetAllAOCs(ctx context.Context) ([]domain.AOCResponse, error) {
	aocs, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		out = append(out, *aoc.ToResponse())
	}
	return out, nil
}

func (s *AOCService) UpdateAOC(ctx context.Context, aoc *domain.AOC) error {
	existing, err := s.repo.GetByID(ctx, aoc.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("aoc not found")
	}
	existing.AccountTypeID = aoc.AccountTypeID
	existing.AccountTaxID = aoc.AccountTaxID
	existing.Code = aoc.Code
	existing.Name = aoc.Name
	existing.Description = aoc.Description
	return s.repo.Update(ctx, existing)
}

func (s *AOCService) DeleteAOC(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.Delete(ctx, ids)
}

func (s *AOCService) BulkUpdateTax(ctx context.Context, ids []uuid.UUID, accountTaxID int) error {
	if len(ids) == 0 {
		return nil
	}
	tax, err := s.repo.GetAccountTaxByID(ctx, accountTaxID)
	if err != nil || tax == nil {
		return errors.New("account tax not found")
	}
	return s.repo.BulkUpdateAccountTax(ctx, ids, accountTaxID)
}

func (s *AOCService) GetAOCType(ctx context.Context) ([]domain.AccountType, error) {
	return s.repo.GetAllAccountTypes(ctx)
}

func (s *AOCService) GetAccountTax(ctx context.Context) ([]domain.AccountTax, error) {
	return s.repo.GetAllAccountTax(ctx)
}
