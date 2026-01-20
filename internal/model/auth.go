package model

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
