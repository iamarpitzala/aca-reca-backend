package user

import (
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type UserClinic struct {
	ID       string `db:"uc_id" json:"id"`
	UserID   string `db:"uc_user_id" json:"userId"`
	ClinicID string `db:"uc_clinic_id" json:"clinicId"`
	Role     string `db:"uc_role" json:"role"`

	Name        string  `db:"uc_clinic_name" json:"name"`
	ABNNumber   string  `db:"uc_clinic_abn_number" json:"abnNumber"`
	Address     string  `db:"uc_clinic_address" json:"address"`
	City        string  `db:"uc_clinic_city" json:"city"`
	State       string  `db:"uc_clinic_state" json:"state"`
	Postcode    *string `db:"uc_clinic_postcode" json:"postcode"`
	Phone       *string `db:"uc_clinic_phone" json:"phone"`
	Email       *string `db:"uc_clinic_email" json:"email"`
	Website     *string `db:"uc_clinic_website" json:"website"`
	LogoURL     *string `db:"uc_clinic_logo_url" json:"logoURL"`
	Description *string `db:"uc_clinic_description" json:"description"`
	ShareType   string  `db:"uc_clinic_share_type" json:"shareType"`
	ClinicShare int     `db:"uc_clinic_clinic_share" json:"clinicShare"`
	OwnerShare  int     `db:"uc_clinic_owner_share" json:"ownerShare"`

	CreatedAt time.Time  `db:"uc_created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"uc_updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"uc_deleted_at" json:"deletedAt"`
}

type UserClinicWithClinic struct {
	UserClinic
	Clinic clinic.Clinic `json:"clinic"`
}

type UserClinicWithUser struct {
	UserClinic
	User User `json:"user"`
}
