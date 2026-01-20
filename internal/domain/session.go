package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       uuid.UUID  `db:"user_id" json:"user_id"`
	RefreshToken string     `db:"refresh_token" json:"-"`
	UserAgent    string     `db:"user_agent" json:"user_agent"`
	IPAddress    string     `db:"ip_address" json:"ip_address"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	ExpiresAt    time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SessionData struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"session_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
