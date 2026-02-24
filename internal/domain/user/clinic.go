package user

import (
	"time"
)

type UserClinic struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	ClinicID  string    `db:"clinic_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// type UserClinicWithClinic struct {
// 	UserClinic
// 	Clinic clinic.Clinic `json:"clinic"`
// }

type UserClinicWithUser struct {
	UserClinic
	User User `json:"user"`
}

type UserClinicWithClinic struct {
	UC_ID            string     `db:"uc_id"`
	UC_UserID        string     `db:"uc_user_id"`
	UC_ClinicID      string     `db:"uc_clinic_id"`
	UC_Role          string     `db:"uc_role"`
	UC_CreatedAt     time.Time  `db:"uc_created_at"`
	UC_UpdatedAt     time.Time  `db:"uc_updated_at"`
	C_ID             string     `db:"c_id"`
	C_Name           string     `db:"c_name"`
	C_ABNNumber      string     `db:"c_abn_number"`
	C_Address        string     `db:"c_address"`
	C_City           string     `db:"c_city"`
	C_State          string     `db:"c_state"`
	C_Postcode       *string    `db:"c_postcode"`
	C_Phone          *string    `db:"c_phone"`
	C_Email          *string    `db:"c_email"`
	C_Website        *string    `db:"c_website"`
	C_LogoURL        *string    `db:"c_logo_url"`
	C_Description    *string    `db:"c_description"`
	C_ShareType      string     `db:"c_share_type"`
	C_ClinicShare    int        `db:"c_clinic_share"`
	C_OwnerShare     int        `db:"c_owner_share"`
	C_MethodType     string     `db:"c_method_type"`
	C_IsActive       bool       `db:"c_is_active"`
	C_WithHoldingTax bool       `db:"c_with_holding_tax"`
	C_CreatedAt      time.Time  `db:"c_created_at"`
	C_UpdatedAt      time.Time  `db:"c_updated_at"`
	C_DeletedAt      *time.Time `db:"c_deleted_at"`
}
