package form

import (
	"time"
)

// Section is a section of a form version (tbl_custom_form_section).
type Section struct {
	ID             int        `db:"id"`
	FormVersionID  int        `db:"form_version_id"`
	SectionTypeID  int        `db:"section_type_id"`
	Name           string     `db:"name"`
	Description    *string    `db:"description"`
	SectionOrder   int        `db:"section_order"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// SectionRequest is used for create/update.
type SectionRequest struct {
	FormVersionID string  `json:"formVersionId" binding:"required"`
	SectionTypeID int     `json:"sectionTypeId" binding:"required"`
	Name          string  `json:"name" binding:"required,max=255"`
	Description   *string `json:"description"`
	SectionOrder  int     `json:"sectionOrder" binding:"min=0"`
}

// SectionResponse is for API return.
type SectionResponse struct {
	ID            int     `json:"id"`
	FormVersionID int     `json:"formVersionId"`
	SectionTypeID int     `json:"sectionTypeId"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	SectionOrder  int     `json:"sectionOrder"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (s *Section) ToSectionResponse() *SectionResponse {
	return &SectionResponse{
		ID:            s.ID,
		FormVersionID: s.FormVersionID,
		SectionTypeID: s.SectionTypeID,
		Name:          s.Name,
		Description:   s.Description,
		SectionOrder:  s.SectionOrder,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

// SectionTypeOption is a static option for the Build form step (Step 1).
// Category is INCOME or EXPENSES; order defines display order (top to bottom).
type SectionTypeOption struct {
	ID          string `json:"id"`          // collection, cost, service_and_facility, other_cost
	Category    string `json:"category"`    // INCOME, EXPENSES
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}
