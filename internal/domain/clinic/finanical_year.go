package clinic

import (
	"time"
)

type FinancialYearRequest struct {
	ID        int        `json:"id"`
	ClinicID  string     `json:"clinic_id"`
	FYLabel   string     `json:"fy_label"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	IsCurrent bool       `json:"is_current"`
	IsClosed  bool       `json:"is_closed"`
	ClosedAt  *time.Time `json:"closed_at"`
	ClosedBy  *string    `json:"closed_by"`
}

func (f *FinancialYearRequest) ToFinancialYear() *FinancialYear {
	return &FinancialYear{
		ID:        f.ID,
		ClinicID:  f.ClinicID,
		FYLabel:   f.FYLabel,
		StartDate: f.StartDate,
		EndDate:   f.EndDate,
		IsCurrent: f.IsCurrent,
		IsClosed:  f.IsClosed,
		ClosedAt:  f.ClosedAt,
		ClosedBy:  f.ClosedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

type FinancialYear struct {
	ID        int        `db:"id"`
	ClinicID  string     `db:"clinic_id"`
	FYLabel   string     `db:"fy_label"`
	StartDate time.Time  `db:"start_date"`
	EndDate   time.Time  `db:"end_date"`
	IsCurrent bool       `db:"is_current"`
	IsClosed  bool       `db:"is_closed"`
	ClosedAt  *time.Time `db:"closed_at"`
	ClosedBy  *string    `db:"closed_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (f *FinancialYear) ToFinancialYearResponse() *FinancialYearResponse {
	return &FinancialYearResponse{
		ID:        f.ID,
		ClinicID:  f.ClinicID,
		FYLabel:   f.FYLabel,
		StartDate: f.StartDate,
		EndDate:   f.EndDate,
		IsCurrent: f.IsCurrent,
		IsClosed:  f.IsClosed,
		ClosedAt:  f.ClosedAt,
		ClosedBy:  f.ClosedBy,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		DeletedAt: f.DeletedAt,
	}
}

type FinancialYearResponse struct {
	ID        int        `json:"id"`
	ClinicID  string     `json:"clinic_id"`
	FYLabel   string     `json:"fy_label"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	IsCurrent bool       `json:"is_current"`
	IsClosed  bool       `json:"is_closed"`
	ClosedAt  *time.Time `json:"closed_at"`
	ClosedBy  *string    `json:"closed_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
