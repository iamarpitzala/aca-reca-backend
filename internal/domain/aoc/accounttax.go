package aoc

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type AccountTaxRequest struct {
	ID          *string    `json:"id" validate:"omitempty,required"`
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Rate        float64    `json:"rate" validate:"required,min=0"`
	Description *string    `json:"description" validate:"omitempty,max=255"`
	CreatedAt   time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt   *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt   *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (a *AccountTaxRequest) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type AccountTax struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Rate        float64    `db:"rate"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (a *AccountTax) ToAccountTaxDB(accountTax *AccountTaxRequest) {
	a.ID = uuid.MustParse(*accountTax.ID)
	a.Name = accountTax.Name
	a.Rate = accountTax.Rate
	a.Description = *accountTax.Description

	a.CreatedAt = accountTax.CreatedAt
	a.UpdatedAt = lo.FromPtr(accountTax.UpdatedAt)
	a.DeletedAt = accountTax.DeletedAt
}

type AccountTaxResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Rate        float64    `json:"rate"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

func (a *AccountTax) ToAccountTaxResponse() *AccountTaxResponse {
	return &AccountTaxResponse{
		ID:          a.ID.String(),
		Name:        a.Name,
		Rate:        a.Rate,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		DeletedAt:   a.DeletedAt,
	}
}
