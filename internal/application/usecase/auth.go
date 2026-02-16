package usecase

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
	"golang.org/x/crypto/bcrypt"
)

// generatePlaceholderABN returns a unique 11-digit ABN for default clinics (derived from userID).
func generatePlaceholderABN(userID uuid.UUID) string {
	h := fnv.New64a()
	_, _ = h.Write(userID[:])
	n := h.Sum64() % 10000000000
	return "1" + fmt.Sprintf("%010d", n)
}

type AuthService struct {
	userRepo     port.UserRepository
	sessionRepo  port.SessionRepository
	token        port.TokenProvider
	clinicUC     *ClinicService
	userClinicUC *UserClinicService
}

func NewAuthService(userRepo port.UserRepository, sessionRepo port.SessionRepository, token port.TokenProvider, clinicUC *ClinicService, userClinicUC *UserClinicService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		token:        token,
		clinicUC:     clinicUC,
		userClinicUC: userClinicUC,
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.New("failed to check if email exists")
	}
	if exists {
		return nil, errors.New("email already in use")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	now := time.Now()
	user := domain.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}
	if err := s.ensureDefaultClinicForUser(ctx, user.ID, user.FirstName); err != nil {
		return nil, fmt.Errorf("registration succeeded but default clinic setup failed; please log in and create a clinic manually: %w", err)
	}
	sessionID := uuid.New()
	tokenPair, err := s.token.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.sessionRepo.Create(ctx, &session); err != nil {
		return nil, err
	}
	user.Password = ""
	return &domain.AuthResponse{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == uuid.Nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	sessionID := uuid.New()
	tokenPair, err := s.token.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}
	now := time.Now()
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    now.Add(s.token.RefreshTokenTTL()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.sessionRepo.Create(ctx, &session); err != nil {
		return nil, err
	}
	user.Password = ""
	return &domain.AuthResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	claims, err := s.token.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, errors.New("session expired")
	}
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	newTokenPair, err := s.token.GenerateTokenPair(user.ID, user.Email, session.ID)
	if err != nil {
		return nil, errors.New("failed to generate new tokens")
	}
	session.RefreshToken = newTokenPair.RefreshToken
	session.ExpiresAt = time.Now().Add(s.token.RefreshTokenTTL())
	session.UpdatedAt = time.Now()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}
	return &domain.AuthResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		TokenType:    newTokenPair.TokenType,
		ExpiresIn:    newTokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

// EnsureDefaultClinicForUser creates a default clinic for the user and links them as owner.
// Idempotent: if the user already has at least one clinic, no-op.
// Used after registration and OAuth sign-up.
func (s *AuthService) EnsureDefaultClinicForUser(ctx context.Context, userID uuid.UUID, displayName string) error {
	return s.ensureDefaultClinicForUser(ctx, userID, displayName)
}

func (s *AuthService) ensureDefaultClinicForUser(ctx context.Context, userID uuid.UUID, displayName string) error {
	clinics, err := s.userClinicUC.GetUserClinics(ctx, userID)
	if err != nil {
		return err
	}
	if len(clinics) > 0 {
		return nil
	}
	name := "My Clinic"
	if displayName != "" {
		name = displayName + "'s Clinic"
	}
	clinic := &domain.Clinic{
		Name:        name,
		ABNNumber:   generatePlaceholderABN(userID),
		Address:     "To be updated",
		City:        "To be updated",
		State:       domain.StateNSW,
		ShareType:   domain.ShareTypePercentage,
		MethodType:  domain.MethodTypeNet,
		ClinicShare: 50,
		OwnerShare:  50,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.clinicUC.CreateClinic(ctx, clinic); err != nil {
		return err
	}
	_, err = s.userClinicUC.AssociateUserWithClinic(ctx, userID, clinic.ID, util.RoleOwner)
	if err != nil {
		_ = s.clinicUC.DeleteClinic(ctx, clinic.ID)
		return err
	}
	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, userID uuid.UUID, firstName, lastName, phone string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *AuthService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	return s.sessionRepo.GetUserSessions(ctx, userID)
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Revoke(ctx, sessionID)
}

func (s *AuthService) OAuthLogin(ctx context.Context, userID uuid.UUID) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	sessionID := uuid.New()
	tokenPair, err := s.token.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}
	now := time.Now()
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    now.Add(s.token.RefreshTokenTTL()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.sessionRepo.Create(ctx, &session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	user.Password = ""
	return &domain.AuthResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
