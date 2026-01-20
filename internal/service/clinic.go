package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type ClinicService struct {
	db *sqlx.DB
}

func NewClinicService(db *sqlx.DB) *ClinicService {
	return &ClinicService{
		db: db,
	}
}

func (cs *ClinicService) CreateClinic(ctx context.Context, clinic *domain.Clinic) error {
	return repository.CreateClinic(ctx, cs.db, clinic)
}

func (cs *ClinicService) GetClinicByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error) {
	return repository.GetClinicByID(ctx, cs.db, id)
}

func (cs *ClinicService) UpdateClinic(ctx context.Context, clinic *domain.Clinic) error {
	return repository.UpdateClinic(ctx, cs.db, clinic)
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
