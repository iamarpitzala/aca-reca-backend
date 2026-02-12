package domain

import (
	"time"

	"github.com/google/uuid"
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

type ShareType string

const (
	ShareTypePercentage ShareType = "PERCENTAGE"
	ShareTypeFixed      ShareType = "FIXED"
)

type MethodType string

const (
	MethodTypeNet   MethodType = "NET"
	MethodTypeGross MethodType = "GROSS"
)

type ClinicRequest struct {
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

	MethodType  MethodType `json:"methodType" validate:"omitempty,required,oneof=NET GROSS"`
	ShareType   ShareType  `json:"shareType" validate:"omitempty,required,oneof=PERCENTAGE FIXED"`
	ClinicShare int        `json:"clinicShare" validate:"omitempty,required,min=0,max=100"`
	OwnerShare  int        `json:"ownerShare" validate:"omitempty,required,min=0,max=100"`

	WithHoldingTax bool `json:"withHoldingTax" validate:"omitempty,required,boolean"`
	IsActive       bool `json:"isActive" validate:"omitempty,required,boolean"`

	CreatedAt time.Time `json:"createdAt" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (c *ClinicRequest) ToClinic() *Clinic {
	return &Clinic{
		Name:           c.Name,
		ABNNumber:      c.ABNNumber,
		Address:        c.Address,
		City:           c.City,
		State:          c.State,
		Postcode:       c.Postcode,
		Phone:          c.Phone,
		Email:          c.Email,
		Website:        c.Website,
		LogoURL:        c.LogoURL,
		Description:    c.Description,
		ShareType:      c.ShareType,
		MethodType:     c.MethodType,
		ClinicShare:    c.ClinicShare,
		OwnerShare:     c.OwnerShare,
		IsActive:       c.IsActive,
		WithHoldingTax: c.WithHoldingTax,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		DeletedAt:      c.DeletedAt,
	}
}

type Clinic struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	ABNNumber   string    `db:"abn_number"`
	Address     string    `db:"address"`
	City        string    `db:"city"`
	State       string    `db:"state"`
	Postcode    *string   `db:"postcode"`
	Phone       *string   `db:"phone"`
	Email       *string   `db:"email"`
	Website     *string   `db:"website"`
	LogoURL     *string   `db:"logo_url"`
	Description *string   `db:"description"`

	ShareType   ShareType  `db:"share_type"`
	MethodType  MethodType `db:"method_type"`
	ClinicShare int        `db:"clinic_share"`
	OwnerShare  int        `db:"owner_share"`

	IsActive       bool `db:"is_active"`
	WithHoldingTax bool `db:"with_holding_tax"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type ClinicResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	ABNNumber      string     `json:"abnNumber"`
	Address        string     `json:"address"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	Postcode       *string    `json:"postcode"`
	Phone          *string    `json:"phone"`
	Email          *string    `json:"email"`
	Website        *string    `json:"website"`
	LogoURL        *string    `json:"logoURL"`
	Description    *string    `json:"description"`
	ShareType      ShareType  `json:"shareType"`
	MethodType     MethodType `json:"methodType"`
	ClinicShare    int        `json:"clinicShare"`
	OwnerShare     int        `json:"ownerShare"`
	IsActive       bool       `json:"isActive"`
	WithHoldingTax bool       `json:"withHoldingTax"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
}

func (c *Clinic) ToClinicResponse() *ClinicResponse {
	return &ClinicResponse{
		ID:             c.ID,
		Name:           c.Name,
		ABNNumber:      c.ABNNumber,
		Address:        c.Address,
		City:           c.City,
		State:          c.State,
		Postcode:       c.Postcode,
		Phone:          c.Phone,
		Email:          c.Email,
		Website:        c.Website,
		LogoURL:        c.LogoURL,
		Description:    c.Description,
		ShareType:      c.ShareType,
		MethodType:     c.MethodType,
		ClinicShare:    c.ClinicShare,
		OwnerShare:     c.OwnerShare,
		IsActive:       c.IsActive,
		WithHoldingTax: c.WithHoldingTax,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		DeletedAt:      &c.DeletedAt,
	}
}

// UpdateClinicRequest supports partial updates - nil/omitted fields are not updated
type UpdateClinicRequest struct {
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

	ShareType      ShareType  `json:"shareType"`
	MethodType     MethodType `json:"methodType"`
	ClinicShare    *int       `json:"clinicShare"`
	OwnerShare     *int       `json:"ownerShare"`
	IsActive       *bool      `json:"isActive"`
	WithHoldingTax *bool      `json:"withHoldingTax"`
}

func (u *UpdateClinicRequest) ToClinic() *Clinic {
	return &Clinic{
		Name:           *u.Name,
		ABNNumber:      *u.ABNNumber,
		Address:        *u.Address,
		City:           *u.City,
		State:          *u.State,
		Postcode:       u.Postcode,
		Phone:          u.Phone,
		Email:          u.Email,
		Website:        u.Website,
		LogoURL:        u.LogoURL,
		Description:    u.Description,
		ShareType:      u.ShareType,
		MethodType:     u.MethodType,
		ClinicShare:    *u.ClinicShare,
		OwnerShare:     *u.OwnerShare,
		IsActive:       *u.IsActive,
		WithHoldingTax: *u.WithHoldingTax,
		UpdatedAt:      time.Now(),
		DeletedAt:      time.Time{},
	}
}
