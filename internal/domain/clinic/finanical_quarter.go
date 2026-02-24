package clinic

import "time"

// FinancialQuarter is a quarter of a master financial year (tbl_financial_quarter).
type FinancialQuarter struct {
	ID              int       `db:"id"`
	FinancialYearID int       `db:"financial_year_id"`
	QuarterNumber   int       `db:"quarter_number"`
	Name            string    `db:"name"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// FinancialQuarterRequest for API.
type FinancialQuarterRequest struct {
	FinancialYearID int       `json:"financial_year_id"`
	Name            string    `json:"name"`
	QuarterNumber   int       `json:"quarter_number"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

// FinancialQuarterResponse for API.
type FinancialQuarterResponse struct {
	ID              int       `json:"id"`
	FinancialYearID int       `json:"financial_year_id"`
	Name            string    `json:"name"`
	QuarterNumber   int       `json:"quarter_number"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (f *FinancialQuarter) ToResponse() *FinancialQuarterResponse {
	return &FinancialQuarterResponse{
		ID:              f.ID,
		FinancialYearID: f.FinancialYearID,
		Name:            f.Name,
		QuarterNumber:   f.QuarterNumber,
		StartDate:       f.StartDate,
		EndDate:         f.EndDate,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}
