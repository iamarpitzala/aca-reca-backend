package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type UserClinicService struct {
	ucRepo   port.UserClinicRepository
	clinicRepo port.ClinicRepository
	userRepo port.UserRepository
}

func NewUserClinicService(ucRepo port.UserClinicRepository, clinicRepo port.ClinicRepository, userRepo port.UserRepository) *UserClinicService {
	return &UserClinicService{
		ucRepo:    ucRepo,
		clinicRepo: clinicRepo,
		userRepo:  userRepo,
	}
}

func (s *UserClinicService) AssociateUserWithClinic(ctx context.Context, userID, clinicID uuid.UUID, role string) (*domain.UserClinic, error) {
	existing, err := s.ucRepo.GetByUserAndClinic(ctx, userID, clinicID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user is already associated with this clinic")
	}
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, errors.New("clinic not found")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if role == "" {
		role = "owner"
	}
	now := time.Now()
	uc := &domain.UserClinic{
		ID:        uuid.New(),
		UserID:    userID,
		ClinicID:  clinicID,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.ucRepo.Create(ctx, uc); err != nil {
		return nil, err
	}
	return uc, nil
}

func (s *UserClinicService) GetUserClinics(ctx context.Context, userID uuid.UUID) ([]domain.UserClinicWithClinic, error) {
	return s.ucRepo.GetUserClinics(ctx, userID)
}

func (s *UserClinicService) GetClinicUsers(ctx context.Context, clinicID uuid.UUID) ([]domain.UserClinicWithUser, error) {
	return s.ucRepo.GetClinicUsers(ctx, clinicID)
}

func (s *UserClinicService) RemoveUserFromClinic(ctx context.Context, id uuid.UUID) error {
	_, err := s.ucRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user-clinic association not found")
	}
	return s.ucRepo.Delete(ctx, id)
}

func (s *UserClinicService) UserHasAccessToClinic(ctx context.Context, userID, clinicID uuid.UUID) (bool, error) {
	uc, err := s.ucRepo.GetByUserAndClinic(ctx, userID, clinicID)
	if err != nil {
		return false, err
	}
	return uc != nil, nil
}
