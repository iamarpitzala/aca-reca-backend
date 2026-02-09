package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AOCRequest struct {
	AccountTypeID int     `json:"accountTypeId" validate:"required"`
	AccountTaxID  int     `json:"accountTaxId" validate:"required"`
	Code          string  `json:"code" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   *string `json:"description" validate:"omitempty,max=255"`
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

func (a *AOCRequest) ToRepo() *AOC {
	now := time.Now()
	aoc := &AOC{
		ID:            uuid.New(),
		AccountTypeID: a.AccountTypeID,
		AccountTaxID:  a.AccountTaxID,
		Code:          a.Code,
		Name:          a.Name,
		Description:   a.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}

	return aoc
}

type AccountType struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type AccountTax struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Rate        float64    `db:"rate"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type AOC struct {
	ID            uuid.UUID  `db:"id"`
	AccountTypeID int        `db:"account_type_id"`
	AccountTaxID  int        `db:"account_tax_id"`
	Code          string     `db:"code"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func (a *AOC) ToResponse() *AOCResponse {
	return &AOCResponse{
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

type AOCResponse struct {
	ID            uuid.UUID  `json:"id"`
	AccountTypeID int        `json:"accountTypeId"`
	AccountTaxID  int        `json:"accountTaxId"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

// BulkDeleteAOCRequest is the body for bulk delete.
type BulkDeleteAOCRequest struct {
	IDs []string `json:"ids" validate:"omitempty,required"`
}

// BulkUpdateTaxRequest is the body for bulk tax change.
type BulkUpdateTaxRequest struct {
	IDs          []string `json:"ids" validate:"required"`
	AccountTaxID int      `json:"accountTaxId" validate:"required"`
}

// BulkArchiveAOCRequest is the body for bulk archive (same as delete: set deleted_at).
type BulkArchiveAOCRequest struct {
	IDs []string `json:"ids" validate:"omitempty,required"`
}
