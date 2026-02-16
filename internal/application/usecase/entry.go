package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/jmoiron/sqlx"
)

type EntryService struct {
	repo           port.EntryRepository
	clinicRepo     port.ClinicRepository
	customFormRepo port.CustomFormRepository
	db             *sqlx.DB
}

func NewEntryService(entryRepo port.EntryRepository, clinicRepo port.ClinicRepository, customFormRepo port.CustomFormRepository) *EntryService {
	return &EntryService{
		repo:           entryRepo,
		clinicRepo:     clinicRepo,
		customFormRepo: customFormRepo,
	}
}

func (s *EntryService) AddEntry(ctx context.Context, formId uuid.UUID, entry domain.CommonEntry) (*domain.CalculationResultNet, error) {
	customForm, err := s.customFormRepo.GetByID(ctx, formId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Failed to fetch form for entry")
	}

	clinic, err := s.clinicRepo.GetByID(ctx, customForm.ClinicID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Failed to fetch clinic for entry")
	}

	field := service.CommonCalculation(customForm, clinic, entry)
	fmt.Println(field.NetAmount)
	fmt.Println(field.TotalPaybleToDentist)

	return field, nil
}
