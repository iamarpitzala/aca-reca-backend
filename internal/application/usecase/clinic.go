package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

var ErrDuplicateABN = errors.New("a clinic with this ABN already exists")

type ClinicService struct {
	repo port.ClinicRepository
}

func NewClinicService(repo port.ClinicRepository) *ClinicService {
	return &ClinicService{repo: repo}
}

func (s *ClinicService) CreateClinic(ctx context.Context, clinic *domain.Clinic) error {
	if err := ValidateState(clinic.State); err != nil {
		return err
	}
	if err := ValidateABN(clinic.ABNNumber); err != nil {
		return err
	}
	clinic.ID = uuid.New()
	if err := s.repo.Create(ctx, clinic); err != nil {
		if strings.Contains(err.Error(), "ux_clinic_abn_active") || strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateABN
		}
		return err
	}
	return nil
}

func (s *ClinicService) GetClinicByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ClinicService) UpdateClinic(ctx context.Context, clinic *domain.Clinic) error {
	if err := ValidateState(clinic.State); err != nil {
		return err
	}
	if err := ValidateABN(clinic.ABNNumber); err != nil {
		return err
	}
	return s.repo.Update(ctx, clinic)
}

func (s *ClinicService) UpdateClinicPartial(ctx context.Context, id uuid.UUID, req *domain.UpdateClinicRequest) (*domain.Clinic, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("clinic not found")
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
	if err := s.repo.Update(ctx, existing); err != nil {
		if strings.Contains(err.Error(), "ux_clinic_abn_active") || strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrDuplicateABN
		}
		return nil, err
	}
	return existing, nil
}

func (s *ClinicService) DeleteClinic(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ClinicService) GetAllClinics(ctx context.Context) ([]domain.Clinic, error) {
	return s.repo.List(ctx)
}

func (s *ClinicService) GetClinicByABNNumber(ctx context.Context, abnNumber string) (*domain.Clinic, error) {
	return s.repo.GetByABN(ctx, abnNumber)
}
