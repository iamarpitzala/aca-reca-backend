package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"-"` // Never return password in JSON
	FirstName       string     `db:"first_name" json:"first_name"`
	LastName        string     `db:"last_name" json:"last_name"`
	Phone           string     `db:"phone" json:"phone"`
	AvatarURL       string     `db:"avatar_url" json:"avatar_url"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	IsEmailVerified bool       `db:"is_email_verified" json:"is_email_verified"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"-"`
}
