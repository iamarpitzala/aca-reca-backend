package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AOSService struct {
	db *sqlx.DB
}

func NewAOSService(db *sqlx.DB) *AOSService {
	return &AOSService{
		db: db,
	}
}

func (as *AOSService) CreateAOC(ctx context.Context, aoc *domain.AOCRequest) error {
	aocRepo := aoc.ToRepo()

	existing, err := repository.GetAOCByCode(ctx, as.db, aocRepo.Code)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existing != nil {
		return errors.New("aoc with this code already exists")
	}

	accountType, err := repository.GetAccountTypeByID(ctx, as.db, aocRepo.AccountTypeID)
	if err != nil {
		return err
	}
	if accountType == nil {
		return errors.New("account type not found")
	}

	accountTax, err := repository.GetAccountTaxByID(ctx, as.db, aocRepo.AccountTaxID)
	if err != nil {
		return err
	}
	if accountTax == nil {
		return errors.New("account tax not found")
	}

	return repository.CreateAOC(ctx, as.db, aocRepo)
}

func (as *AOSService) GetAOCByID(ctx context.Context, id uuid.UUID) (*domain.AOCResponse, error) {
	aoc, err := repository.GetAOCByID(ctx, as.db, id)
	if err != nil {
		return nil, err
	}
	return aoc.ToResponse(), nil
}

func (as *AOSService) GetAOCByCode(ctx context.Context, code string) (*domain.AOCResponse, error) {
	aoc, err := repository.GetAOCByCode(ctx, as.db, code)
	if err != nil {
		return nil, err
	}
	return aoc.ToResponse(), nil
}

func (as *AOSService) GetAOCByAccountTypeID(ctx context.Context, accountTypeID int, sortBy, sortOrder string) ([]domain.AOCResponse, error) {
	aocs, err := repository.GetAOCByAccountTypeIDSorted(ctx, as.db, accountTypeID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		responses = append(responses, *aoc.ToResponse())
	}
	return responses, nil
}

func (as *AOSService) GetAOCsByAccountType(ctx context.Context, sortBy, sortOrder string) ([]domain.AOCResponse, error) {
	aocs, err := repository.GetAOCsByAccountTypeSorted(ctx, as.db, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		responses = append(responses, *aoc.ToResponse())
	}
	return responses, nil
}

func (as *AOSService) GetAOCByAccountTaxID(ctx context.Context, accountTaxID int) ([]domain.AOCResponse, error) {
	aocs, err := repository.GetAOCByAccountTaxID(ctx, as.db, accountTaxID)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		responses = append(responses, *aoc.ToResponse())
	}
	return responses, nil
}

func (as *AOSService) GetAllAOCs(ctx context.Context) ([]domain.AOCResponse, error) {
	aocs, err := repository.GetAllAOCs(ctx, as.db)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.AOCResponse, 0, len(aocs))
	for _, aoc := range aocs {
		responses = append(responses, *aoc.ToResponse())
	}
	return responses, nil
}

func (as *AOSService) UpdateAOC(ctx context.Context, aoc *domain.AOC) error {
	aocRepo, err := repository.GetAOCByID(ctx, as.db, aoc.ID)
	if err != nil {
		return err
	}
	aocRepo.AccountTypeID = aoc.AccountTypeID
	aocRepo.AccountTaxID = aoc.AccountTaxID
	aocRepo.Code = aoc.Code
	aocRepo.Name = aoc.Name
	aocRepo.Description = aoc.Description
	return repository.UpdateAOC(ctx, as.db, aocRepo)
}

func (as *AOSService) DeleteAOC(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return repository.DeleteAOC(ctx, as.db, ids)
}

func (as *AOSService) BulkUpdateTax(ctx context.Context, ids []uuid.UUID, accountTaxID int) error {
	if len(ids) == 0 {
		return nil
	}
	tax, err := repository.GetAccountTaxByID(ctx, as.db, accountTaxID)
	if err != nil || tax == nil {
		return errors.New("account tax not found")
	}
	return repository.BulkUpdateAccountTax(ctx, as.db, ids, accountTaxID)
}

func (as *AOSService) GetAOCType(ctx context.Context) ([]domain.AccountType, error) {
	accountTypes, err := repository.GetAllAOCType(ctx, as.db)
	return accountTypes, err
}

func (as *AOSService) GetAccountTax(ctx context.Context) ([]domain.AccountTax, error) {
	taxes, err := repository.GetAllAccountTax(ctx, as.db)
	return taxes, err
}
