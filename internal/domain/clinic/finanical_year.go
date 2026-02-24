package clinic

import (
	"time"
)

// FinancialYear is the global/master financial year (tbl_financial_year).
type FinancialYear struct {
	ID        int        `db:"id"`
	FYLabel   string     `db:"fy_label"`
	StartDate time.Time  `db:"start_date"`
	EndDate   time.Time  `db:"end_date"`
	IsActive  bool       `db:"is_active"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// SetTimestamps sets created_at and updated_at for Create.
func (f *FinancialYear) SetTimestamps() {
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now
}

// FinancialYearRequest for API (create master FY or link clinic to FY).
type FinancialYearRequest struct {
	ID              int       `json:"id"`
	FinancialYearID int       `json:"financialYearId"` // when set, create link (tbl_clinic_financial_year)
	FYLabel         string    `json:"fyLabel"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	IsActive        bool      `json:"isActive"`
}

// ToFinancialYear maps request to master FinancialYear (for create/update master FY).
func (r *FinancialYearRequest) ToFinancialYear() *FinancialYear {
	return &FinancialYear{
		ID:        r.ID,
		FYLabel:   r.FYLabel,
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
		IsActive:  r.IsActive,
	}
}

// FinancialYearResponse for API.
type FinancialYearResponse struct {
	ID        int        `json:"id"`
	FYLabel   string     `json:"fyLabel"`
	StartDate time.Time  `json:"startDate"`
	EndDate   time.Time  `json:"endDate"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (f *FinancialYear) ToResponse() *FinancialYearResponse {
	return &FinancialYearResponse{
		ID:        f.ID,
		FYLabel:   f.FYLabel,
		StartDate: f.StartDate,
		EndDate:   f.EndDate,
		IsActive:  f.IsActive,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		DeletedAt: f.DeletedAt,
	}
}
