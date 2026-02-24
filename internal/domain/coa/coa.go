package coa

import (
	"time"

	"github.com/google/uuid"
)

type COARequest struct {
	ID            *string `json:"id" validate:"omitempty,required"`
	AccountTypeID int     `json:"accountTypeId" validate:"required"`
	AccountTaxID  int     `json:"accountTaxId" validate:"required"`
	Code          string  `json:"code" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   *string `json:"description" validate:"omitempty,max=255"`
}

func (a *COARequest) ToRepo() *COA {
	now := time.Now()
	coa := &COA{
		ID:            uuid.New().String(),
		AccountTypeID: a.AccountTypeID,
		AccountTaxID:  a.AccountTaxID,
		Code:          a.Code,
		Name:          a.Name,
		Description:   a.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}

	return coa
}

type AccountTypeCOA struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type AccountTaxCOA struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Rate        float64    `db:"rate"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type COA struct {
	ID            string     `db:"id"`
	OwnerUserID   string     `db:"owner_user_id"`
	AccountTypeID int        `db:"account_type_id"`
	AccountTaxID  int        `db:"account_tax_id"`
	Code          string     `db:"code"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func (a *COA) ToResponse() *COAResponse {
	return &COAResponse{
		ID:            a.ID,
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

type COAResponse struct {
	ID            string     `json:"id"`
	AccountTypeID int        `json:"accountTypeId"`
	AccountTaxID  int        `json:"accountTaxId"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

// BulkDeleteCOARequest is the body for bulk delete.
type BulkDeleteCOARequest struct {
	IDs []string `json:"ids" validate:"omitempty,required"`
}

// BulkUpdateCOATaxRequest is the body for bulk tax change.
type BulkUpdateCOATaxRequest struct {
	IDs          []string `json:"ids" validate:"required"`
	AccountTaxID int      `json:"accountTaxId" validate:"required"`
}

// BulkArchiveCOARequest is the body for bulk archive (same as delete: set deleted_at).
type BulkArchiveCOARequest struct {
	IDs []string `json:"ids" validate:"omitempty,required"`
}
