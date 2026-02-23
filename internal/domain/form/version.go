package form

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type FormVersionRequest struct {
	ID        *string    `json:"id" validate:"omitempty,required"`
	FormID    string     `json:"formId" validate:"required"`
	Version   int        `json:"version" validate:"required,min=1"`
	IsActive  bool       `json:"isActive" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	CreatedBy string     `json:"createdBy" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

type FormVersion struct {
	ID        string     `db:"id"`
	FormID    string     `db:"form_id"`
	Version   int        `db:"version"`
	IsActive  bool       `db:"is_active"`
	CreatedAt time.Time  `db:"created_at"`
	CreatedBy string     `db:"created_by"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (f *FormVersion) ToFormVersionDB(formVersion *FormVersionRequest) {
	if formVersion.ID != nil {
		f.ID = *formVersion.ID
	} else {
		f.ID = uuid.New().String()
	}
	f.FormID = formVersion.FormID
	f.Version = formVersion.Version
	f.IsActive = formVersion.IsActive
	f.CreatedAt = formVersion.CreatedAt
	f.CreatedBy = formVersion.CreatedBy
}

type FormVersionResponse struct {
	ID       string `json:"id"`
	FormID   string `json:"formId"`
	Version  int    `json:"version"`
	IsActive bool   `json:"isActive"`

	CreatedBy *string    `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (f *FormVersion) ToFormVersionResponse() *FormVersionResponse {
	return &FormVersionResponse{
		ID:        f.ID,
		FormID:    f.FormID,
		Version:   f.Version,
		IsActive:  f.IsActive,
		CreatedBy: lo.ToPtr(f.CreatedBy),
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		DeletedAt: f.DeletedAt,
	}
}
