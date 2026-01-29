package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

var ErrDuplicateABN = errors.New("a clinic with this ABN already exists")

type ClinicService struct {
	db *sqlx.DB
}

func NewClinicService(db *sqlx.DB) *ClinicService {
	return &ClinicService{
		db: db,
	}
}

func (cs *ClinicService) CreateClinic(ctx context.Context, clinic *domain.Clinic) error {
	// Validate state
	if err := ValidateState(clinic.State); err != nil {
		return err
	}

	// Validate ABN
	if err := ValidateABN(clinic.ABNNumber); err != nil {
		return err
	}

	clinic.ID = uuid.New()
	err := repository.CreateClinic(ctx, cs.db, clinic)
	if err != nil {
		if strings.Contains(err.Error(), "ux_clinic_abn_active") || strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateABN
		}
		return err
	}
	return nil
}

func (cs *ClinicService) GetClinicByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error) {
	return repository.GetClinicByID(ctx, cs.db, id)
}

func (cs *ClinicService) UpdateClinic(ctx context.Context, clinic *domain.Clinic) error {
	// Validate state
	if err := ValidateState(clinic.State); err != nil {
		return err
	}

	// Validate ABN
	if err := ValidateABN(clinic.ABNNumber); err != nil {
		return err
	}

	return repository.UpdateClinic(ctx, cs.db, clinic)
}

// UpdateClinicPartial merges partial updates into existing clinic and saves
func (cs *ClinicService) UpdateClinicPartial(ctx context.Context, id uuid.UUID, req *domain.UpdateClinicRequest) (*domain.Clinic, error) {
	existing, err := repository.GetClinicByID(ctx, cs.db, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.ABNNumber != nil {
		existing.ABNNumber = *req.ABNNumber
		if err := ValidateABN(existing.ABNNumber); err != nil {
			return nil, err
		}
	}
	if req.Address != nil {
		existing.Address = *req.Address
	}
	if req.City != nil {
		existing.City = *req.City
	}
	if req.State != nil {
		if err := ValidateState(*req.State); err != nil {
			return nil, err
		}
		existing.State = *req.State
	}
	if req.Postcode != nil {
		existing.Postcode = req.Postcode
	}
	if req.Phone != nil {
		existing.Phone = req.Phone
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.Website != nil {
		existing.Website = req.Website
	}
	if req.LogoURL != nil {
		existing.LogoURL = req.LogoURL
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.ClinicShare != nil {
		existing.ClinicShare = *req.ClinicShare
	}
	if req.OwnerShare != nil {
		existing.OwnerShare = *req.OwnerShare
	}
	if err := repository.UpdateClinic(ctx, cs.db, existing); err != nil {
		if strings.Contains(err.Error(), "ux_clinic_abn_active") || strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrDuplicateABN
		}
		return nil, err
	}
	return existing, nil
}

func (cs *ClinicService) DeleteClinic(ctx context.Context, id uuid.UUID) error {
	return repository.DeleteClinic(ctx, cs.db, id)
}

func (cs *ClinicService) GetAllClinics(ctx context.Context) ([]domain.Clinic, error) {
	return repository.GetAllClinics(ctx, cs.db)
}

func (cs *ClinicService) GetClinicByABNNumber(ctx context.Context, abnNumber string) (*domain.Clinic, error) {
	return repository.GetClinicByABNNumber(ctx, cs.db, abnNumber)
}
