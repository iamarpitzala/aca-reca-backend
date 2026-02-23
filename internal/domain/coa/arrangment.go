package coa

import (
	"time"
)

// Arrangement represents a row in tbl_arrangement (method, name, percentage, amount, user_id, etc.)
type Arrangement struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	Method      string     `db:"method"`
	Name        string     `db:"name"`
	Percentage  float64    `db:"percentage"`
	Amount      *float64   `db:"amount"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// ArrangementRequest is the payload for create/update (validate with go-playground/validator).
type ArrangementRequest struct {
	Method      string   `json:"method" validate:"required,oneof=GROSS NET"`
	Name        string   `json:"name" validate:"required,max=255"`
	Percentage  float64  `json:"percentage" validate:"min=0,max=100"`
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
}

// ToArrangement builds an Arrangement from request (create: id, userID and timestamps set in usecase/repo).
func (r *ArrangementRequest) ToArrangement(userID string) *Arrangement {
	return &Arrangement{
		UserID:      userID,
		Method:      r.Method,
		Name:        r.Name,
		Percentage:  r.Percentage,
		Amount:      r.Amount,
		Description: r.Description,
	}
}

// ArrangementResponse is the API response shape (camelCase JSON).
type ArrangementResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Method      string     `json:"method"`
	Name        string     `json:"name"`
	Percentage  float64    `json:"percentage"`
	Amount      *float64   `json:"amount"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

// ToResponse maps domain Arrangement to response DTO.
func (a *Arrangement) ToResponse() *ArrangementResponse {
	return &ArrangementResponse{
		ID:          a.ID,
		UserID:      a.UserID,
		Method:      a.Method,
		Name:        a.Name,
		Percentage:  a.Percentage,
		Amount:      a.Amount,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		DeletedAt:   a.DeletedAt,
	}
}

// Arrangment is kept for backward compatibility; prefer Arrangement.
type Arrangment = Arrangement
