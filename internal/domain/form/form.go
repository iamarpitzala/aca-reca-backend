package form

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// FormRequest is used for API interactions, validation, and client communication
type FormRequest struct {
	ID          *string `json:"id" validate:"omitempty,required"`
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description *string `json:"description" validate:"omitempty,min=3,max=225"`
	ClinicID    string  `json:"clinicId" validate:"required"`

	// Version is handled by FormVersion, not tbl_custom_form
	VersionID int `json:"version" validate:"required,min=1"`

	Status           string `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	AccountingMethod string `json:"accountingMethod" validate:"required,oneof=NET GROSS"`

	CreatedBy string     `json:"createdBy" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`

	PublishedAt *time.Time `json:"publishedAt" validate:"omitempty,required_with=DeletedAt"`
}

func (f *FormRequest) IsPublished() bool {
	return f.PublishedAt != nil
}

// Form is the DB representation for tbl_custom_form
type Form struct {
	ID               string     `db:"id"`
	ClinicID         string     `db:"clinic_id"`
	Name             string     `db:"name"`
	Description      *string    `db:"description"`
	Status           string     `db:"status"`
	VersionID        int        `db:"version_id"`
	AccountingMethod string     `db:"calculation_method"`
	CreatedBy        string     `db:"created_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	PublishedAt      *time.Time `db:"published_at"`
}

func (f *Form) ToFormDB(form *FormRequest) {
	if form.ID != nil && *form.ID != "" {
		f.ID = *form.ID
	} else {
		f.ID = uuid.New().String()
	}
	f.ClinicID = form.ClinicID
	f.Name = form.Name
	f.Description = form.Description
	f.AccountingMethod = form.AccountingMethod
	f.Status = form.Status
	f.CreatedBy = form.CreatedBy
	f.CreatedAt = form.CreatedAt
	f.UpdatedAt = lo.FromPtr(form.UpdatedAt)
	f.DeletedAt = form.DeletedAt
	f.PublishedAt = form.PublishedAt
	f.VersionID = form.VersionID
}

// FormResponse is for API return
type FormResponse struct {
	ID          string  `json:"id"`
	ClinicID    string  `json:"clinicId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`

	VersionID int    `json:"versionId"`
	Status    string `json:"status"`

	AccountingMethod string `json:"accountingMethod"`

	CreatedBy   string     `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	PublishedAt *time.Time `json:"publishedAt"`
}

func (f *Form) ToFormResponse() *FormResponse {
	return &FormResponse{
		ID:               f.ID,
		ClinicID:         f.ClinicID,
		Name:             f.Name,
		Description:      f.Description,
		Status:           f.Status,
		AccountingMethod: f.AccountingMethod,
		VersionID:        f.VersionID,
		CreatedBy:        f.CreatedBy,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        lo.ToPtr(f.UpdatedAt),
		DeletedAt:        f.DeletedAt,
		PublishedAt:      f.PublishedAt,
	}
}
