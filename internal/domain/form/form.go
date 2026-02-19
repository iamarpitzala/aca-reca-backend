package form

import (
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/samber/lo"
)

type FormStatus string

const (
	FormStatusDraft     FormStatus = util.FormStatusDraft
	FormStatusPublished FormStatus = util.FormStatusPublished
	FormStatusArchived  FormStatus = util.FormStatusArchived
)

type TaxTreatment string

const (
	TaxTreatmentInclusive TaxTreatment = util.GSTTypeInclusive
	TaxTreatmentExclusive TaxTreatment = util.GSTTypeExclusive
	TaxTreatmentManual    TaxTreatment = util.GSTTypeManual
	TaxTreatmentNone      TaxTreatment = util.GSTTypeNone
)

type FormRequest struct {
	ID          *string `json:"id" validate:"omitempty,required"`
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=3,max=255"`
	ClinicID    string  `json:"clinicId" validate:"required"`

	Version      int          `json:"version" validate:"required,min=1"`
	Status       FormStatus   `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	TaxTreatment TaxTreatment `json:"taxTreatment" validate:"required,oneof=INCLUSIVE EXCLUSIVE MANUAL NONE"`

	CreatedBy *uuid.UUID `json:"createdBy" validate:"required_with=CreatedAt"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`

	PublishedAt *time.Time `json:"publishedAt" validate:"omitempty,required_with=DeletedAt"`
}

func (f *FormRequest) IsPublished() bool {
	return f.PublishedAt != nil
}

func (f *FormRequest) IncrementVersion() {
	f.Version++
}

type Form struct {
	ID           uuid.UUID    `db:"id"`
	ClinicID     uuid.UUID    `db:"clinic_id"`
	Name         string       `db:"name"`
	Description  string       `db:"description"`
	TaxTreatment TaxTreatment `db:"tax_treatment"`
	Version      int          `db:"version"`
	Status       FormStatus   `db:"status"`
	CreatedBy    *uuid.UUID   `db:"created_by"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    *time.Time   `db:"deleted_at"`
	PublishedAt  *time.Time   `db:"published_at"`
}

func (f *Form) ToFormDB(form *FormRequest) {
	f.ID = uuid.MustParse(*form.ID)
	f.ClinicID = uuid.MustParse(form.ClinicID)
	f.Name = form.Name
	f.Description = form.Description
	f.TaxTreatment = form.TaxTreatment
	f.Version = form.Version
	f.Status = form.Status
	f.CreatedBy = form.CreatedBy
	f.CreatedAt = form.CreatedAt
	f.UpdatedAt = lo.FromPtr(form.UpdatedAt)
	f.DeletedAt = form.DeletedAt
	f.PublishedAt = form.PublishedAt
}

type FormResponse struct {
	ID           string       `json:"id"`
	ClinicID     string       `json:"clinicId"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	TaxTreatment TaxTreatment `json:"taxTreatment"`

	Version int        `json:"version"`
	Status  FormStatus `json:"status"`

	CreatedBy   *uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	PublishedAt *time.Time `json:"publishedAt"`
}

func (f *Form) ToFormResponse() *FormResponse {

	return &FormResponse{
		ID:           f.ID.String(),
		ClinicID:     f.ClinicID.String(),
		Name:         f.Name,
		Description:  f.Description,
		TaxTreatment: f.TaxTreatment,
		Version:      f.Version,
		Status:       f.Status,
		CreatedBy:    f.CreatedBy,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    lo.ToPtr(f.UpdatedAt),
		DeletedAt:    f.DeletedAt,
		PublishedAt:  f.PublishedAt,
	}
}
