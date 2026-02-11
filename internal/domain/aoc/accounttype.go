package aoc

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type AccountTypeRequest struct {
	ID          *string    `json:"id" validate:"omitempty,required"`
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description *string    `json:"description" validate:"omitempty,max=255"`
	CreatedAt   time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt   *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt   *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (a *AccountTypeRequest) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type AccountType struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (a *AccountType) ToAccountTypeDB(accountType *AccountTypeRequest) {
	a.ID = uuid.MustParse(*accountType.ID)
	a.Name = accountType.Name
	a.Description = *accountType.Description
	a.CreatedAt = accountType.CreatedAt
	a.UpdatedAt = lo.FromPtr(accountType.UpdatedAt)
	a.DeletedAt = accountType.DeletedAt
}

type AccountTypeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

func (a *AccountType) ToAccountTypeResponse() *AccountTypeResponse {
	return &AccountTypeResponse{
		ID:          a.ID.String(),
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		DeletedAt:   a.DeletedAt,
	}
}
