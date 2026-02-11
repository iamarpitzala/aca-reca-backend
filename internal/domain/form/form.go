package form

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type FormStatus string

const (
	FormStatusDraft     FormStatus = "DRAFT"
	FormStatusPublished FormStatus = "PUBLISHED"
	FormStatusArchived  FormStatus = "ARCHIVED"
)

type FormRequest struct {
	ID                *string           `json:"id" validate:"omitempty,required"`
	Name              string            `json:"name" validate:"required,min=3,max=255"`
	Description       string            `json:"description" validate:"omitempty,min=3,max=255"`
	CalculationMethod CalculationMethod `json:"calculationMethod" validate:"required,oneof=NET GROSS"`
	Fields            []FieldRequest    `json:"fields" validate:"omitempty,required"`

	Version int        `json:"version" validate:"required,min=1"`
	Status  FormStatus `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`

	CreatedBy *uuid.UUID `json:"createdBy" validate:"required_with=CreatedAt"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`

	PublishedAt *time.Time `json:"publishedAt" validate:"omitempty,required_with=DeletedAt"`
}

func (f *FormRequest) Validate() error {
	for _, field := range f.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (f *FormRequest) IsPublished() bool {
	return f.PublishedAt != nil
}

func (f *FormRequest) IncrementVersion() {
	f.Version++
}

type Form struct {
	ID                uuid.UUID         `db:"id"`
	Name              string            `db:"name"`
	Description       string            `db:"description"`
	CalculationMethod CalculationMethod `db:"calculation_method"`
	Fields            []Field           `db:"fields"`
	Version           int               `db:"version"`
	Status            FormStatus        `db:"status"`
	CreatedBy         *uuid.UUID        `db:"created_by"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`
	DeletedAt         *time.Time        `db:"deleted_at"`
	PublishedAt       *time.Time        `db:"published_at"`
}

func (f *Form) ToFormDB(form *FormRequest) {
	fields := make([]Field, len(form.Fields))
	for i, field := range form.Fields {
		fields[i].ToFieldDB(&field)
	}

	f.ID = uuid.MustParse(*form.ID)
	f.Name = form.Name
	f.Description = form.Description
	f.CalculationMethod = form.CalculationMethod
	f.Fields = fields
	f.Version = form.Version
	f.Status = form.Status
	f.CreatedBy = form.CreatedBy
	f.CreatedAt = form.CreatedAt
	f.UpdatedAt = lo.FromPtr(form.UpdatedAt)
	f.DeletedAt = form.DeletedAt
	f.PublishedAt = form.PublishedAt
}

type FormResponse struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	CalculationMethod CalculationMethod `json:"calculationMethod"`
	Fields            []FieldResponse   `json:"fields"`

	Version int        `json:"version"`
	Status  FormStatus `json:"status"`

	CreatedBy   *uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	PublishedAt *time.Time `json:"publishedAt"`
}

func (f *Form) ToFormResponse() *FormResponse {
	fields := make([]FieldResponse, len(f.Fields))
	for i, field := range f.Fields {
		fields[i] = *field.ToFieldResponse()
	}
	return &FormResponse{
		ID:                f.ID.String(),
		Name:              f.Name,
		Description:       f.Description,
		CalculationMethod: f.CalculationMethod,
		Fields:            fields,
		Version:           f.Version,
		Status:            f.Status,
		CreatedBy:         f.CreatedBy,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         lo.ToPtr(f.UpdatedAt),
		DeletedAt:         f.DeletedAt,
		PublishedAt:       f.PublishedAt,
	}
}
