package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserClinic struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"userId"`
	ClinicID  uuid.UUID  `db:"clinic_id" json:"clinicId"`
	Role      string     `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type UserClinicWithClinic struct {
	UserClinic
	Clinic Clinic `json:"clinic"`
}

type UserClinicWithUser struct {
	UserClinic
	User User `json:"user"`
}
