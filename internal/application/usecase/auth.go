package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    port.UserRepository
	sessionRepo port.SessionRepository
	token       port.TokenProvider
}

func NewAuthService(userRepo port.UserRepository, sessionRepo port.SessionRepository, token port.TokenProvider) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		token:       token,
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
