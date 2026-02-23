package clinic

import "time"

type FinancialQuarterRequest struct {
	FinancialYearID int        `json:"financial_year_id"`
	Name            string     `json:"name"`
	QuarterNumber   int        `json:"quarter_number"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	IsClosed        bool       `json:"is_closed"`
	ClosedAt        *time.Time `json:"closed_at"`
	ClosedBy        *string    `json:"closed_by"`
}

func (f *FinancialQuarterRequest) ToFinancialQuarter() *FinancialQuarter {
	return &FinancialQuarter{
		FinancialYearID: f.FinancialYearID,
		Name:            f.Name,
		QuarterNumber:   f.QuarterNumber,
		StartDate:       f.StartDate,
		EndDate:         f.EndDate,
		IsClosed:        f.IsClosed,
		ClosedAt:        f.ClosedAt,
		ClosedBy:        f.ClosedBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}
}

type FinancialQuarter struct {
	ID              int        `db:"id"`
	FinancialYearID int        `db:"financial_year_id"`
	Name            string     `db:"name"`
	QuarterNumber   int        `db:"quarter_number"`
	StartDate       time.Time  `db:"start_date"`
	EndDate         time.Time  `db:"end_date"`
	IsClosed        bool       `db:"is_closed"`
	ClosedAt        *time.Time `db:"closed_at"`
	ClosedBy        *string    `db:"closed_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

func (f *FinancialQuarter) ToFinancialQuarterResponse() *FinancialQuarterResponse {
	return &FinancialQuarterResponse{
		ID:              f.ID,
		FinancialYearID: f.FinancialYearID,
		Name:            f.Name,
		QuarterNumber:   f.QuarterNumber,
		StartDate:       f.StartDate,
		EndDate:         f.EndDate,
		IsClosed:        f.IsClosed,
		ClosedAt:        f.ClosedAt,
		ClosedBy:        f.ClosedBy,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
		DeletedAt:       f.DeletedAt,
	}
}

type FinancialQuarterResponse struct {
	ID              int        `json:"id"`
	FinancialYearID int        `json:"financial_year_id"`
	Name            string     `json:"name"`
	QuarterNumber   int        `json:"quarter_number"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	IsClosed        bool       `json:"is_closed"`
	ClosedAt        *time.Time `json:"closed_at"`
	ClosedBy        *string    `json:"closed_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
