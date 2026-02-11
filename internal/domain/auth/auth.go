package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderPassword  AuthProvider = "PASSWORD"
	AuthProviderGoogle    AuthProvider = "GOOGLE"
	AuthProviderMicrosoft AuthProvider = "MICROSOFT"
)

func (a AuthProvider) String() string {
	return string(a)
}

func (a AuthProvider) IsValid() bool {
	return a == AuthProviderPassword || a == AuthProviderGoogle || a == AuthProviderMicrosoft
}

type AuthIdentity struct {
	ID             uuid.UUID    `db:"id"`
	UserID         uuid.UUID    `db:"user_id"`
	Provider       AuthProvider `db:"provider"`
	ProviderUserID *string      `db:"provider_user_id"`
	Email          *string      `db:"email"`
	PasswordHash   *string      `db:"password_hash"`
	EmailVerified  bool         `db:"email_verified"`
	IsActive       bool         `db:"is_active"`
	AvatarURL      *string      `db:"avatar_url"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      *time.Time   `db:"updated_at"`
	DeletedAt      *time.Time   `db:"deleted_at"`
}

func (a *AuthIdentity) ToAuthIdentityDB(authIdentity *AuthIdentityRequest) {
	a.ID = uuid.MustParse(*authIdentity.ID)
	a.UserID = uuid.MustParse(*authIdentity.UserID)
	a.Provider = authIdentity.Provider
	a.ProviderUserID = authIdentity.ProviderUserID
	a.Email = authIdentity.Email
	a.PasswordHash = authIdentity.PasswordHash
	a.EmailVerified = authIdentity.EmailVerified
	a.IsActive = authIdentity.IsActive
	a.AvatarURL = authIdentity.AvatarURL
	a.CreatedAt = authIdentity.CreatedAt
	a.UpdatedAt = authIdentity.UpdatedAt
	a.DeletedAt = authIdentity.DeletedAt
}

type AuthIdentityRequest struct {
	ID             *string      `json:"id" validate:"omitempty,required"`
	UserID         *string      `json:"userId" validate:"omitempty,required"`
	Provider       AuthProvider `json:"provider" validate:"required,oneof=PASSWORD GOOGLE MICROSOFT"`
	ProviderUserID *string      `json:"providerUserId" validate:"omitempty,required"`
	Email          *string      `json:"email" validate:"omitempty,required,email"`
	PasswordHash   *string      `json:"passwordHash" validate:"omitempty,required"`
	EmailVerified  bool         `json:"emailVerified" validate:"omitempty,required"`
	IsActive       bool         `json:"isActive" validate:"omitempty,required,boolean"`
	AvatarURL      *string      `json:"avatarURL" validate:"omitempty,required,url"`

	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (a *AuthIdentityRequest) Validate() error {
	if a.Provider.IsValid() {
		return fmt.Errorf("provider is required")
	}
	if a.ProviderUserID == nil && a.Provider != AuthProviderPassword {
		return fmt.Errorf("provider user id is required for provider %s", a.Provider)
	}
	if a.Email == nil && a.Provider != AuthProviderPassword {
		return fmt.Errorf("email is required for provider %s", a.Provider)
	}
	if a.PasswordHash == nil && a.Provider == AuthProviderPassword {
		return fmt.Errorf("password hash is required for provider %s", a.Provider)
	}
	if a.EmailVerified {
		return fmt.Errorf("email verified is required for provider %s", a.Provider)
	}
	if a.CreatedAt.IsZero() {
		return fmt.Errorf("created at is required")
	}
	if a.UpdatedAt == nil && a.CreatedAt.IsZero() {
		return fmt.Errorf("updated at is required")
	}
	if a.DeletedAt == nil && a.UpdatedAt.IsZero() {
		return fmt.Errorf("deleted at is required")
	}
	return nil
}

type AuthIdentityResponse struct {
	ID             string       `json:"id"`
	UserID         uuid.UUID    `json:"userId"`
	Provider       AuthProvider `json:"provider"`
	ProviderUserID *string      `json:"providerUserId"`
	Email          *string      `json:"email"`
	PasswordHash   *string      `json:"passwordHash"`
	EmailVerified  bool         `json:"emailVerified"`
	IsActive       bool         `json:"isActive"`
	AvatarURL      *string      `json:"avatarURL"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      *time.Time   `json:"updatedAt"`
	DeletedAt      *time.Time   `json:"deletedAt"`
}

func (a *AuthIdentity) ToAuthIdentityResponse() *AuthIdentityResponse {
	return &AuthIdentityResponse{
		ID:             a.ID.String(),
		UserID:         a.UserID,
		Provider:       a.Provider,
		ProviderUserID: a.ProviderUserID,
		Email:          a.Email,
		PasswordHash:   a.PasswordHash,
		EmailVerified:  a.EmailVerified,
		IsActive:       a.IsActive,
		AvatarURL:      a.AvatarURL,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		DeletedAt:      a.DeletedAt,
	}
}
