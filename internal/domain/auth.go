package domain

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"userId"`
	Provider       string     `db:"provider" json:"provider"`               // google, microsoft, etc.
	ProviderUserID string     `db:"provider_user_id" json:"providerUserId"` // External provider's user ID
	ProviderEmail  string     `db:"provider_email" json:"providerEmail"`
	AccessToken    string     `db:"access_token" json:"accessToken"`   // Encrypted in production
	RefreshToken   string     `db:"refresh_token" json:"refreshToken"` // Encrypted in production
	TokenExpiresAt *time.Time `db:"token_expires_at" json:"tokenExpiresAt"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}
