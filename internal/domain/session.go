package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       uuid.UUID  `db:"user_id" json:"userId"`
	RefreshToken string     `db:"refresh_token" json:"-"`
	UserAgent    string     `db:"user_agent" json:"userAgent"`
	IPAddress    string     `db:"ip_address" json:"ipAddress"`
	ExpiresAt    time.Time  `db:"expires_at" json:"expiresAt"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deletedAt"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SessionData struct {
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"sessionId"`
	UserAgent string    `json:"userAgent"`
	IPAddress string    `json:"ipAddress"`
	CreatedAt time.Time `json:"createdAt"`
}

type TokenClaims struct {
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"sessionId"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}
