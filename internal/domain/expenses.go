package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExpenseType struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ClinicID    uuid.UUID  `db:"clinic_id" json:"clinicId"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy   uuid.UUID  `db:"created_by" json:"createdBy"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	DeletedBy   uuid.UUID  `db:"deleted_by" json:"deletedBy"`
}

type ExpenseCategory struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ClinicID    uuid.UUID  `db:"clinic_id" json:"clinicId"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy   uuid.UUID  `db:"created_by" json:"createdBy"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	DeletedBy   uuid.UUID  `db:"deleted_by" json:"deletedBy"`
}

type ExpenseCategoryType struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	ClinicID   uuid.UUID  `db:"clinic_id" json:"clinicId"`
	TypeID     uuid.UUID  `db:"type_id" json:"typeId"`
	CategoryID uuid.UUID  `db:"category_id" json:"categoryId"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy  uuid.UUID  `db:"created_by" json:"createdBy"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	DeletedBy  uuid.UUID  `db:"deleted_by" json:"deletedBy"`
}

type ExpenseEntry struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ClinicID       uuid.UUID  `db:"clinic_id" json:"clinicId"`
	CategoryID     uuid.UUID  `db:"category_id" json:"categoryId"`
	TypeID         uuid.UUID  `db:"type_id" json:"typeId"`
	Amount         float64    `db:"amount" json:"amount"`
	GSTRate        *float64   `db:"gst_rate" json:"gstRate"`
	IsGSTInclusive *bool      `db:"is_gst_inclusive" json:"isGSTInclusive"`
	ExpenseDate    time.Time  `db:"expense_date" json:"expenseDate"`
	SupplierName   string     `db:"supplier_name" json:"supplierName"`
	Notes          string     `db:"notes" json:"notes"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy      uuid.UUID  `db:"created_by" json:"createdBy"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	DeletedBy      uuid.UUID  `db:"deleted_by" json:"deletedBy"`
}
