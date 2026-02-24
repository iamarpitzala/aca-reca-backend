package clinic

import (
	"time"
)

// ClinicFinancialYear links a clinic to a master financial year (tbl_clinic_financial_year).
type ClinicFinancialYear struct {
	ID                int        `db:"id"`
	ClinicID          string     `db:"clinic_id"`
	FinancialYearID   int        `db:"financial_year_id"`
	IsCurrent         bool       `db:"is_current"`
	IsClosed          bool       `db:"is_closed"`
	ClosedAt          *time.Time `db:"closed_at"`
	ClosedBy          *string    `db:"closed_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}
