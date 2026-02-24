package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/user"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

type UserClinicService struct {
	ucRepo     port.UserClinicRepository
	clinicRepo port.ClinicRepository
	userRepo   port.UserRepository
}

func NewUserClinicService(ucRepo port.UserClinicRepository, clinicRepo port.ClinicRepository, userRepo port.UserRepository) *UserClinicService {
	return &UserClinicService{
		ucRepo:     ucRepo,
		clinicRepo: clinicRepo,
		userRepo:   userRepo,
	}
}

func (s *UserClinicService) AssociateUserWithClinic(ctx context.Context, userID, clinicID string, role string) (*user.UserClinic, error) {
	existing, err := s.ucRepo.GetByUserAndClinic(ctx, userID, clinicID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user is already associated with this clinic")
	}
	if role == "" {
		role = util.RoleOwner
	}
	now := time.Now()
	uc := &user.UserClinic{
		ID:        uuid.New().String(),
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

func (s *UserClinicService) GetUserClinics(ctx context.Context, userID string) ([]user.UserClinicWithClinic, error) {
	return s.ucRepo.GetUserClinics(ctx, userID)
}

func (s *UserClinicService) RemoveUserFromClinic(ctx context.Context, id string) error {
	_, err := s.ucRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user-clinic association not found")
	}
	return s.ucRepo.Delete(ctx, id)
}

func (s *UserClinicService) UserHasAccessToClinic(ctx context.Context, userID, clinicID string) (bool, error) {
	userClinics, err := s.ucRepo.GetUserClinics(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, uc := range userClinics {
		if uc.UC_ClinicID == clinicID {
			return true, nil
		}
	}
	return false, nil
}

func (s *UserClinicService) UserRoleInClinic(ctx context.Context, userID, clinicID string) (string, error) {
	userClinics, err := s.ucRepo.GetUserClinics(ctx, userID)
	if err != nil {
		return "", err
	}
	for _, uc := range userClinics {
		if uc.UC_ClinicID == clinicID {
			return uc.UC_Role, nil
		}
	}
	return "", errors.New("user not associated with this clinic")
}
