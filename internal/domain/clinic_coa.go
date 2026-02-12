package domain

import (
	"time"

	"github.com/google/uuid"
)

// ClinicCOA is the junction entity linking a clinic to a chart-of-accounts (AOC) entry.
type ClinicCOA struct {
	ID        uuid.UUID  `db:"id"`
	ClinicID  uuid.UUID  `db:"clinic_id"`
	COAID     uuid.UUID  `db:"coa_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// ClinicCOARequest is the body for associating an AOC with a clinic.
type ClinicCOARequest struct {
	COAID uuid.UUID `json:"coaId" binding:"required"`
}

// ClinicCOAResponse is the API response for a clinic–AOC association.
type ClinicCOAResponse struct {
	ID        uuid.UUID  `json:"id"`
	ClinicID  uuid.UUID  `json:"clinicId"`
	COAID     uuid.UUID  `json:"coaId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// ToResponse converts ClinicCOA to ClinicCOAResponse.
func (c *ClinicCOA) ToResponse() *ClinicCOAResponse {
	return &ClinicCOAResponse{
		ID:        c.ID,
		ClinicID:  c.ClinicID,
		COAID:     c.COAID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}
