package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db           *sqlx.DB
	tokenService *TokenService
	oauthService *OAuthService
}

func NewAuthService(db *sqlx.DB, tokenService *TokenService, oauthService *OAuthService) *AuthService {
	return &AuthService{
		db:           db,
		tokenService: tokenService,
		oauthService: oauthService,
	}
}

func (as *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user already exists
	emailExists, err := repository.EmailExists(ctx, as.db, req.Email)
	if err != nil {
		return nil, errors.New("failed to check if email exists")
	}
	if emailExists {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := domain.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repository.CreateUser(ctx, as.db, &user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	sessionID := uuid.New()
	tokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session in database
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repository.CreateSession(ctx, as.db, &session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user
	user, err := repository.GetUserByEmail(ctx, as.db, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if user.IsActive == false {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate tokens
	sessionID := uuid.New()
	tokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session in database
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repository.CreateSession(ctx, as.db, &session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	// Validate refresh token
	claims, err := as.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find session in database
	session, err := repository.GetSessionByRefreshToken(ctx, as.db, refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, errors.New("session expired")
	}

	// Find user
	user, err := repository.GetUserByID(ctx, as.db, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	newTokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, session.ID)
	if err != nil {
		return nil, errors.New("failed to generate new tokens")
	}

	// Update session with new refresh token
	session.RefreshToken = newTokenPair.RefreshToken
	session.ExpiresAt = time.Now().Add(as.tokenService.refreshTokenTTL)
	session.UpdatedAt = time.Now()

	err = repository.UpdateSession(ctx, as.db, session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		TokenType:    newTokenPair.TokenType,
		ExpiresIn:    newTokenPair.ExpiresIn,
	}, nil
}

func (as *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	// Deactivate session in database
	err := repository.DeleteSession(ctx, as.db, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := repository.GetUserByID(ctx, as.db, userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (as *AuthService) UpdateUser(ctx context.Context, userID uuid.UUID, firstName, lastName, phone string) (*domain.User, error) {
	user, err := repository.GetUserByID(ctx, as.db, userID)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone
	user.UpdatedAt = time.Now()
	err = repository.UpdateUser(ctx, as.db, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (as *AuthService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	sessions, err := repository.GetUserSessions(ctx, as.db, userID)
	if err != nil {
		return sessions, nil
	}

	return sessions, nil
}

func (as *AuthService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	err := repository.RevokeSession(ctx, as.db, sessionID)
	if err != nil {
		return err
	}

	return nil
}

// OAuthLogin creates a session for an OAuth-authenticated user
func (as *AuthService) OAuthLogin(ctx context.Context, userID uuid.UUID) (*domain.AuthResponse, error) {
	// Find user
	user, err := repository.GetUserByID(ctx, as.db, userID)
	if err != nil {
		return nil, err
	}
	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Generate tokens
	sessionID := uuid.New()
	tokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session in database
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := repository.CreateSession(ctx, as.db, &session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Remove password from response
	user.Password = ""

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
