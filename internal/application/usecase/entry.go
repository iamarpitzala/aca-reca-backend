package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/jmoiron/sqlx"
)

type EntryService struct {
	repo port.EntryRepository
	db   *sqlx.DB
}

func NewEntryService(entryRepo port.EntryRepository) *EntryService {
	return &EntryService{
		repo: entryRepo,
	}
}

func (s *EntryService) AddEntry(ctx context.Context, formId uuid.UUID, entry domain.CommonEntry) error {
	customForm, err := repository.GetCustomFormByID(ctx, s.db, formId)
	if err != nil {
		return errors.New("Failed to fetch form for entry")
	}

	clinic, err := repository.GetClinicByID(ctx, s.db, customForm.ClinicID)
	if err != nil {
		return errors.New("Failed to fetch clinic for entry")
	}

	field := service.CommonCalculation(customForm, clinic, entry)
	fmt.Println(field)
	return nil
}
