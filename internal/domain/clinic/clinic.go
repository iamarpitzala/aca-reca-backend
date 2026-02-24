package clinic

import (
	"time"
)

// Australian states/territories
const (
	StateNSW = "NSW"
	StateVIC = "VIC"
	StateQLD = "QLD"
	StateSA  = "SA"
	StateWA  = "WA"
	StateTAS = "TAS"
	StateNT  = "NT"
	StateACT = "ACT"
)

var ValidStates = []string{StateNSW, StateVIC, StateQLD, StateSA, StateWA, StateTAS, StateNT, StateACT}

type ClinicRequest struct {
	UserID      string  `json:"userId" validate:"omitempty,required"`
	Name        string  `json:"name" validate:"omitempty,required"`
	ABNNumber   string  `json:"abnNumber" validate:"required, min=11, max=11"`
	Address     string  `json:"address" validate:"omitempty,required,max=255"`
	City        string  `json:"city" validate:"omitempty,required,max=255"`
	State       string  `json:"state" validate:"omitempty,required,oneof=NSW VIC QLD SA WA TAS NT ACT"`
	Postcode    *string `json:"postcode" validate:"omitempty,required"`
	Phone       *string `json:"phone" validate:"omitempty,required,max=255"`
	Email       *string `json:"email" validate:"omitempty,required,email"`
	Website     *string `json:"website" validate:"omitempty,required"`
	LogoURL     *string `json:"logoURL" validate:"omitempty,required,url"`
	Description *string `json:"description" validate:"omitempty,required,max=255"`

	MethodType  string `json:"methodType" validate:"omitempty,required,oneof=NET GROSS"`
	ShareType   string `json:"shareType" validate:"omitempty,required,oneof=PERCENTAGE FIXED"`
	ClinicShare int    `json:"clinicShare" validate:"omitempty,required,min=0,max=100"`
	OwnerShare  int    `json:"ownerShare" validate:"omitempty,required,min=0,max=100"`

	IsActive bool `json:"isActive" validate:"omitempty,required,boolean"`

	CreatedAt time.Time `json:"createdAt" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (c *ClinicRequest) ToClinic() *Clinic {
	return &Clinic{
		UserID:      c.UserID,
		Name:        c.Name,
		ABNNumber:   c.ABNNumber,
		Address:     c.Address,
		City:        c.City,
		State:       c.State,
		Postcode:    c.Postcode,
		Phone:       c.Phone,
		Email:       c.Email,
		Website:     c.Website,
		LogoURL:     c.LogoURL,
		Description: c.Description,
		ShareType:   c.ShareType,
		MethodType:  c.MethodType,
		ClinicShare: c.ClinicShare,
		OwnerShare:  c.OwnerShare,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
	}
}

type Clinic struct {
	ID          string  `db:"id"`
	UserID      string  `db:"user_id"`
	Name        string  `db:"name"`
	ABNNumber   string  `db:"abn_number"`
	Address     string  `db:"address"`
	City        string  `db:"city"`
	State       string  `db:"state"`
	Postcode    *string `db:"postcode"`
	Phone       *string `db:"phone"`
	Email       *string `db:"email"`
	Website     *string `db:"website"`
	LogoURL     *string `db:"logo_url"`
	Description *string `db:"description"`

	ShareType   string `db:"share_type"`
	MethodType  string `db:"method_type"`
	ClinicShare int    `db:"clinic_share"`
	OwnerShare  int    `db:"owner_share"`

	IsActive bool `db:"is_active"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type ClinicResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Name        string     `json:"name"`
	ABNNumber   string     `json:"abnNumber"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	Postcode    *string    `json:"postcode"`
	Phone       *string    `json:"phone"`
	Email       *string    `json:"email"`
	Website     *string    `json:"website"`
	LogoURL     *string    `json:"logoURL"`
	Description *string    `json:"description"`
	ShareType   string     `json:"shareType"`
	MethodType  string     `json:"methodType"`
	ClinicShare int        `json:"clinicShare"`
	OwnerShare  int        `json:"ownerShare"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

func (c *Clinic) ToClinicResponse() *ClinicResponse {
	return &ClinicResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		Name:        c.Name,
		ABNNumber:   c.ABNNumber,
		Address:     c.Address,
		City:        c.City,
		State:       c.State,
		Postcode:    c.Postcode,
		Phone:       c.Phone,
		Email:       c.Email,
		Website:     c.Website,
		LogoURL:     c.LogoURL,
		Description: c.Description,
		ShareType:   c.ShareType,
		MethodType:  c.MethodType,
		ClinicShare: c.ClinicShare,
		OwnerShare:  c.OwnerShare,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   &c.DeletedAt,
	}
}

// UpdateClinicRequest supports partial updates - nil/omitted fields are not updated
type UpdateClinicRequest struct {
	UserID      *string `json:"userId"`
	Name        *string `json:"name"`
	ABNNumber   *string `json:"abnNumber"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	Postcode    *string `json:"postcode"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Website     *string `json:"website"`
	LogoURL     *string `json:"logoURL"`
	Description *string `json:"description"`

	ShareType   string `json:"shareType"`
	MethodType  string `json:"methodType"`
	ClinicShare *int   `json:"clinicShare"`
	OwnerShare  *int   `json:"ownerShare"`
	IsActive    *bool  `json:"isActive"`
}

func (u *UpdateClinicRequest) ToClinic() *Clinic {
	return &Clinic{
		UserID:      *u.UserID,
		Name:        *u.Name,
		ABNNumber:   *u.ABNNumber,
		Address:     *u.Address,
		City:        *u.City,
		State:       *u.State,
		Postcode:    u.Postcode,
		Phone:       u.Phone,
		Email:       u.Email,
		Website:     u.Website,
		LogoURL:     u.LogoURL,
		Description: u.Description,
		ShareType:   u.ShareType,
		MethodType:  u.MethodType,
		ClinicShare: *u.ClinicShare,
		OwnerShare:  *u.OwnerShare,
		IsActive:    *u.IsActive,
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Time{},
	}
}
