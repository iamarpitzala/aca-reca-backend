package aoc

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountRequest struct {
	ID            *string    `json:"id" validate:"omitempty,required"`
	AccountTypeID int        `json:"accountTypeId" validate:"required"`
	AccountTaxID  int        `json:"accountTaxId" validate:"required"`
	Code          string     `json:"code" validate:"required"`
	Name          string     `json:"name" validate:"required"`
	Description   *string    `json:"description" validate:"omitempty,max=255"`
	CreatedAt     time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt     *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt     *time.Time `json:"deletedAt" validate:"omitempty,required_with=UpdatedAt"`
}

func (a *AccountRequest) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}

	if a.Code == "" {
		return fmt.Errorf("code is required")
	}

	if a.AccountTypeID == 0 {
		return fmt.Errorf("account type id is required")
	}

	if a.AccountTaxID == 0 {
		return fmt.Errorf("account tax id is required")
	}

	return nil
}

type Account struct {
	ID            uuid.UUID  `db:"id"`
	AccountTypeID int        `db:"account_type_id"`
	AccountTaxID  int        `db:"account_tax_id"`
	Code          string     `db:"code"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func (a *Account) ToAccountDB(account *AccountRequest) {
	a.ID = uuid.MustParse(*account.ID)
	a.AccountTypeID = account.AccountTypeID
	a.AccountTaxID = account.AccountTaxID
	a.Code = account.Code
	a.Name = account.Name
	a.Description = account.Description
	a.CreatedAt = account.CreatedAt
	a.UpdatedAt = account.UpdatedAt
	a.DeletedAt = account.DeletedAt
}

type AccountResponse struct {
	ID            string     `json:"id"`
	AccountTypeID int        `json:"accountTypeId"`
	AccountTaxID  int        `json:"accountTaxId"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

func (a *Account) ToAccountResponse() *AccountResponse {
	return &AccountResponse{
		ID:            a.ID.String(),
		AccountTypeID: a.AccountTypeID,
		AccountTaxID:  a.AccountTaxID,
		Code:          a.Code,
		Name:          a.Name,
		Description:   a.Description,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
		DeletedAt:     a.DeletedAt,
	}
}
