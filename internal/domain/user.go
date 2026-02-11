package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	ID              *string    `json:"id" validate:"omitempty,required"`
	Email           string     `json:"email" validate:"omitempty,required,email"`
	Password        string     `json:"password" validate:"omitempty,required,min=8"`
	FirstName       string     `json:"firstName" validate:"omitempty,required"`
	LastName        string     `json:"lastName" validate:"omitempty,required"`
	Phone           string     `json:"phone" validate:"omitempty,required"`
	AvatarURL       *string    `json:"avatarURL" validate:"omitempty,required,url"`
	IsEmailVerified bool       `json:"isEmailVerified" validate:"omitempty,required,boolean"`
	CreatedAt       time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt       *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt       *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

type User struct {
	ID              uuid.UUID  `db:"id"`
	Email           string     `db:"email"`
	Password        string     `db:"password"` // Never return password in JSON
	FirstName       string     `db:"first_name"`
	LastName        string     `db:"last_name"`
	Phone           string     `db:"phone"`
	AvatarURL       string     `db:"avatar_url"`
	IsEmailVerified bool       `db:"is_email_verified"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

func (u *User) ToUserDB(user *UserRequest) {
	u.ID = uuid.MustParse(*user.ID)
	u.Email = user.Email
	u.Password = user.Password
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Phone = user.Phone
	u.AvatarURL = *user.AvatarURL
	u.IsEmailVerified = user.IsEmailVerified
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = *user.UpdatedAt
	u.DeletedAt = user.DeletedAt
}

type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName"`
	Phone           string     `json:"phone"`
	AvatarURL       string     `json:"avatarURL"`
	IsEmailVerified bool       `json:"isEmailVerified"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Phone:           u.Phone,
		AvatarURL:       u.AvatarURL,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		DeletedAt:       u.DeletedAt,
	}
}
