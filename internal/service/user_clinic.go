package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type UserClinicService struct {
	db *sqlx.DB
}

func NewUserClinicService(db *sqlx.DB) *UserClinicService {
	return &UserClinicService{
		db: db,
	}
}

func (ucs *UserClinicService) AssociateUserWithClinic(ctx context.Context, userID, clinicID uuid.UUID, role string) (*domain.UserClinic, error) {
	// Check if association already exists
	existing, err := repository.GetUserClinicByUserAndClinic(ctx, ucs.db, userID, clinicID)
	if err == nil && existing != nil {
		return nil, errors.New("user is already associated with this clinic")
	}

	// Verify clinic exists
	_, err = repository.GetClinicByID(ctx, ucs.db, clinicID)
	if err != nil {
		return nil, errors.New("clinic not found")
	}

	// Verify user exists
	_, err = repository.GetUserByID(ctx, ucs.db, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if role == "" {
		role = "owner"
	}

	userClinic := &domain.UserClinic{
		ID:       uuid.New(),
		UserID:   userID,
		ClinicID: clinicID,
		Role:     role,
	}

	err = repository.CreateUserClinic(ctx, ucs.db, userClinic)
	if err != nil {
		return nil, err
	}

	return userClinic, nil
}

func (ucs *UserClinicService) GetUserClinics(ctx context.Context, userID uuid.UUID) ([]domain.UserClinicWithClinic, error) {
	return repository.GetUserClinics(ctx, ucs.db, userID)
}

func (ucs *UserClinicService) GetClinicUsers(ctx context.Context, clinicID uuid.UUID) ([]domain.UserClinicWithUser, error) {
	return repository.GetClinicUsers(ctx, ucs.db, clinicID)
}

func (ucs *UserClinicService) RemoveUserFromClinic(ctx context.Context, id uuid.UUID) error {
	_, err := repository.GetUserClinicByID(ctx, ucs.db, id)
	if err != nil {
		return errors.New("user-clinic association not found")
	}

	return repository.DeleteUserClinic(ctx, ucs.db, id)
}

func (ucs *UserClinicService) UserHasAccessToClinic(ctx context.Context, userID, clinicID uuid.UUID) (bool, error) {
	userClinic, err := repository.GetUserClinicByUserAndClinic(ctx, ucs.db, userID, clinicID)
	if err != nil {
		return false, err
	}
	return userClinic != nil, nil
}
