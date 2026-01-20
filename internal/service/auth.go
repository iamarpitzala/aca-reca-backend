package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db             *sqlx.DB
	tokenService   *TokenService
	sessionService *SessionService
	oauthService   *OAuthService
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         *model.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
}

func NewAuthService(db *sqlx.DB, tokenService *TokenService, sessionService *SessionService, oauthService *OAuthService) *AuthService {
	return &AuthService{
		db:             db,
		tokenService:   tokenService,
		sessionService: sessionService,
		oauthService:   oauthService,
	}
}

func (as *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser model.User
	err := as.db.GetContext(ctx, &existingUser, "SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL", req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != sql.ErrNoRows {
		return nil, errors.New("failed to check existing user")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := model.User{
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

	query := `INSERT INTO users (id, email, password, first_name, last_name, phone, is_active, created_at, updated_at)
		VALUES (:id, :email, :password, :first_name, :last_name, :phone, :is_active, :created_at, :updated_at)`

	_, err = as.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	sessionID := uuid.New()
	tokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session in database
	session := model.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sessionQuery := `INSERT INTO sessions (id, user_id, refresh_token, is_active, expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :refresh_token, :is_active, :expires_at, :created_at, :updated_at)`

	_, err = as.db.NamedExecContext(ctx, sessionQuery, session)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	// Store session in Redis
	sessionData := &SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		CreatedAt: time.Now(),
	}

	if err := as.sessionService.StoreSession(ctx, sessionID, sessionData); err != nil {
		// Log error but don't fail registration
		_ = err
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *LoginRequest, userAgent string, ipAddress string) (*AuthResponse, error) {
	// Find user
	var user model.User
	err := as.db.GetContext(ctx, &user,
		"SELECT id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL",
		req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("failed to find user")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	sessionID := uuid.New()
	tokenPair, err := as.tokenService.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session in database
	session := model.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sessionQuery := `INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip_address, is_active, expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :refresh_token, :user_agent, :ip_address, :is_active, :expires_at, :created_at, :updated_at)`

	_, err = as.db.NamedExecContext(ctx, sessionQuery, session)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	// Store session in Redis
	sessionData := &SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}

	if err := as.sessionService.StoreSession(ctx, sessionID, sessionData); err != nil {
		// Log error but don't fail login
		_ = err
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := as.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find session in database
	var session model.Session
	err = as.db.GetContext(ctx, &session,
		"SELECT id, user_id, refresh_token, user_agent, ip_address, is_active, expires_at, created_at, updated_at FROM sessions WHERE refresh_token = $1 AND is_active = $2 AND deleted_at IS NULL",
		refreshToken, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found or inactive")
		}
		return nil, errors.New("failed to find session")
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, errors.New("session expired")
	}

	// Find user
	var user model.User
	err = as.db.GetContext(ctx, &user,
		"SELECT id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
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

	_, err = as.db.ExecContext(ctx,
		"UPDATE sessions SET refresh_token = $1, expires_at = $2, updated_at = $3 WHERE id = $4",
		session.RefreshToken, session.ExpiresAt, session.UpdatedAt, session.ID)
	if err != nil {
		return nil, errors.New("failed to update session")
	}

	// Update session in Redis
	sessionData := &SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: session.ID,
		CreatedAt: time.Now(),
	}

	if err := as.sessionService.StoreSession(ctx, session.ID, sessionData); err != nil {
		// Log error but don't fail refresh
		_ = err
	}

	return newTokenPair, nil
}

func (as *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	// Deactivate session in database
	_, err := as.db.ExecContext(ctx, "UPDATE sessions SET is_active = $1, updated_at = $2 WHERE id = $3", false, time.Now(), sessionID)
	if err != nil {
		return errors.New("failed to deactivate session")
	}

	// Delete session from Redis
	if err := as.sessionService.DeleteSession(ctx, sessionID); err != nil {
		// Log error but don't fail logout
		_ = err
	}

	return nil
}

func (as *AuthService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := as.db.Get(&user,
		"SELECT id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user")
	}

	user.Password = ""
	return &user, nil
}

func (as *AuthService) UpdateUser(userID uuid.UUID, firstName, lastName, phone string) (*model.User, error) {
	query := `UPDATE users SET 
		first_name = COALESCE(NULLIF($1, ''), first_name),
		last_name = COALESCE(NULLIF($2, ''), last_name),
		phone = COALESCE(NULLIF($3, ''), phone),
		updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at`

	var user model.User
	err := as.db.Get(&user, query, firstName, lastName, phone, time.Now(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to update user")
	}

	user.Password = ""
	return &user, nil
}

func (as *AuthService) GetUserSessions(userID uuid.UUID) ([]model.Session, error) {
	var sessions []model.Session
	err := as.db.Select(&sessions,
		"SELECT id, user_id, refresh_token, user_agent, ip_address, is_active, expires_at, created_at, updated_at FROM sessions WHERE user_id = $1 AND is_active = $2 AND deleted_at IS NULL",
		userID, true)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (as *AuthService) RevokeSession(sessionID uuid.UUID) error {
	_, err := as.db.Exec("UPDATE sessions SET is_active = $1, updated_at = $2 WHERE id = $3", false, time.Now(), sessionID)
	if err != nil {
		return errors.New("failed to revoke session")
	}

	return nil
}

// OAuthLogin creates a session for an OAuth-authenticated user
func (as *AuthService) OAuthLogin(ctx context.Context, userID uuid.UUID, userAgent string, ipAddress string) (*AuthResponse, error) {
	// Find user
	var user model.User
	err := as.db.GetContext(ctx, &user,
		"SELECT id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user")
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
	session := model.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(as.tokenService.refreshTokenTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sessionQuery := `INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip_address, is_active, expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :refresh_token, :user_agent, :ip_address, :is_active, :expires_at, :created_at, :updated_at)`

	_, err = as.db.NamedExecContext(ctx, sessionQuery, session)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	// Store session in Redis
	sessionData := &SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}

	if err := as.sessionService.StoreSession(ctx, sessionID, sessionData); err != nil {
		// Log error but don't fail login
		_ = err
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
