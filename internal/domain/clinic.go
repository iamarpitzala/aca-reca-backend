package domain

import (
	"time"

	"github.com/google/uuid"
)

type Clinic struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	ABNNumber   string     `db:"abn_number" json:"abnNumber"`
	Address     string     `db:"address" json:"address"`
	City        string     `db:"city" json:"city"`
	State       string     `db:"state" json:"state"`
	Postcode    *string    `db:"postcode" json:"postcode"`
	Phone       *string    `db:"phone" json:"phone"`
	Email       *string    `db:"email" json:"email"`
	Website     *string    `db:"website" json:"website"`
	LogoURL     *string    `db:"logo_url" json:"logoURL"`
	Description *string    `db:"description" json:"description"`
	IsActive    bool       `db:"is_active" json:"isActive"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}
