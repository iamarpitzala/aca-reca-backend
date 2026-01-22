package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"-"` // Never return password in JSON
	FirstName       string     `db:"first_name" json:"firstName"`
	LastName        string     `db:"last_name" json:"lastName"`
	Phone           string     `db:"phone" json:"phone"`
	AvatarURL       string     `db:"avatar_url" json:"avatarURL"`
	IsEmailVerified bool       `db:"is_email_verified" json:"isEmailVerified"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt"`
}
