package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthRequest struct {
	Code        string `json:"code" validate:"required"`
	RedirectURI string `json:"redirectUri" validate:"required"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	IDToken      string
}

type UserInfo struct {
	Sub           string // Google unique user id
	Email         string
	EmailVerified bool
	FirstName     string
	LastName      string
	AvatarURL     string
}

type GoogleAuth struct {
	UserID   uuid.UUID
	UserInfo UserInfo
	Token    Token
}
