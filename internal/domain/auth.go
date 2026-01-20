package domain

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`
	Provider       string     `db:"provider" json:"provider"`                 // google, microsoft, etc.
	ProviderUserID string     `db:"provider_user_id" json:"provider_user_id"` // External provider's user ID
	ProviderEmail  string     `db:"provider_email" json:"provider_email"`
	AccessToken    string     `db:"access_token" json:"-"`  // Encrypted in production
	RefreshToken   string     `db:"refresh_token" json:"-"` // Encrypted in production
	TokenExpiresAt *time.Time `db:"token_expires_at" json:"token_expires_at"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"-"`
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
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
