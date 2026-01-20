package model

import (
	"time"

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
