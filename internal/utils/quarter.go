package utils

import (
	"fmt"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// QuarterStatus represents the status of a quarter
type QuarterStatus string

const (
	QuarterStatusOpen  QuarterStatus = "open"
	QuarterStatusLocked QuarterStatus = "locked"
	QuarterStatusDraft QuarterStatus = "draft"
)

// CalculatedQuarter represents a quarter calculated from financial settings
type CalculatedQuarter struct {
	ID                 string    `json:"id"`
	QuarterIndex       int       `json:"quarterIndex"`
	Label              string    `json:"label"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	Status             QuarterStatus `json:"status"`
	FinancialYearStart time.Time `json:"financialYearStart"`
	FinancialYearEnd   time.Time `json:"financialYearEnd"`
}

// GetFinancialYearStart calculates the financial year start date based on financial year setting
func GetFinancialYearStart(financialYearStart domain.FinancialYearStart, referenceDate time.Time) time.Time {
	year := referenceDate.Year()
	month := int(referenceDate.Month()) // 1-indexed (1=Jan, 7=Jul)

	if financialYearStart == domain.FinancialYearStartJuly {
		// Australian default: July-June
		// If we're in Jan-Jun (months 1-6), the financial year started last July
		// If we're in Jul-Dec (months 7-12), the financial year started this July
		if month < 7 {
			// Jan-Jun: FY started last July
			return time.Date(year-1, 7, 1, 0, 0, 0, 0, referenceDate.Location())
		} else {
			// Jul-Dec: FY started this July
			return time.Date(year, 7, 1, 0, 0, 0, 0, referenceDate.Location())
		}
	} else {
		// Calendar year: January-December
		return time.Date(year, 1, 1, 0, 0, 0, 0, referenceDate.Location())
	}
}

// CalculateQuarters calculates quarters for a given financial year
// Quarters are system-driven and deterministic
func CalculateQuarters(
	financialYearStart domain.FinancialYearStart,
	lockDate *time.Time,
	referenceDate time.Time,
) []CalculatedQuarter {
	fyStart := GetFinancialYearStart(financialYearStart, referenceDate)
	
	quarters := make([]CalculatedQuarter, 4)
	
	for quarterIndex := 0; quarterIndex < 4; quarterIndex++ {
		// Calculate quarter start: financial year start + (quarter_index * 3 months)
		quarterStart := fyStart.AddDate(0, quarterIndex*3, 0)
		
		// Calculate quarter end: quarter start + 3 months - 1 day
		quarterEnd := quarterStart.AddDate(0, 3, -1)
		
		// Financial year end is the day before next FY start
		nextFYStart := fyStart.AddDate(1, 0, 0)
		fyEnd := nextFYStart.AddDate(0, 0, -1)
		
		// Determine status based on lock date
		status := QuarterStatusOpen
		if lockDate != nil {
			if quarterEnd.Before(*lockDate) {
				// Quarter end is before lock date - fully locked
				status = QuarterStatusLocked
			} else if quarterStart.After(*lockDate) {
				// Quarter start is after lock date - fully open
				status = QuarterStatusOpen
			} else {
				// Quarter contains lock date - quarter end is after lock date but start is before/equal
				// This means the quarter spans the lock date, so it's in draft state
				status = QuarterStatusDraft
			}
		}
		
		// Generate deterministic ID: FY{year}-{month}-Q{quarter}
		fyYear := fyStart.Year()
		fyMonth := int(fyStart.Month())
		id := fmt.Sprintf("FY%d-%02d-Q%d", fyYear, fyMonth, quarterIndex+1)
		
		// Format label: "Q{n} ({start_date} – {end_date} {year})"
		startFormatted := quarterStart.Format("2 Jan")
		endFormatted := quarterEnd.Format("2 Jan 2006")
		label := fmt.Sprintf("Q%d (%s – %s)", quarterIndex+1, startFormatted, endFormatted)
		
		quarters[quarterIndex] = CalculatedQuarter{
			ID:                 id,
			QuarterIndex:       quarterIndex,
			Label:              label,
			StartDate:          quarterStart,
			EndDate:            quarterEnd,
			Status:             status,
			FinancialYearStart: fyStart,
			FinancialYearEnd:   fyEnd,
		}
	}
	
	return quarters
}

// GetQuartersForPeriod gets quarters for multiple years (current and previous)
func GetQuartersForPeriod(
	financialYearStart domain.FinancialYearStart,
	lockDate *time.Time,
	yearsBack int,
	yearsForward int,
) []CalculatedQuarter {
	var allQuarters []CalculatedQuarter
	now := time.Now()
	
	// Calculate quarters for previous years
	for yearOffset := -yearsBack; yearOffset <= yearsForward; yearOffset++ {
		referenceDate := now.AddDate(yearOffset, 0, 0)
		quarters := CalculateQuarters(financialYearStart, lockDate, referenceDate)
		allQuarters = append(allQuarters, quarters...)
	}
	
	// Sort by start date (oldest first)
	for i := 0; i < len(allQuarters)-1; i++ {
		for j := i + 1; j < len(allQuarters); j++ {
			if allQuarters[i].StartDate.After(allQuarters[j].StartDate) {
				allQuarters[i], allQuarters[j] = allQuarters[j], allQuarters[i]
			}
		}
	}
	
	return allQuarters
}

// GetQuarterForDate finds the quarter that contains a given date
func GetQuarterForDate(
	date time.Time,
	financialYearStart domain.FinancialYearStart,
	lockDate *time.Time,
) *CalculatedQuarter {
	quarters := GetQuartersForPeriod(financialYearStart, lockDate, 2, 2)
	
	for i := range quarters {
		q := &quarters[i]
		if (date.After(q.StartDate) || date.Equal(q.StartDate)) &&
			(date.Before(q.EndDate) || date.Equal(q.EndDate)) {
			return q
		}
	}
	
	return nil
}

// IsQuarterLocked checks if a quarter is locked
func IsQuarterLocked(quarter *CalculatedQuarter) bool {
	return quarter.Status == QuarterStatusLocked
}

// FormatQuarterWithStatus formats quarter for display with status indicator
func FormatQuarterWithStatus(quarter *CalculatedQuarter) string {
	var statusIcon string
	switch quarter.Status {
	case QuarterStatusOpen:
		statusIcon = "🟢"
	case QuarterStatusLocked:
		statusIcon = "🔒"
	case QuarterStatusDraft:
		statusIcon = "⚠️"
	}
	return fmt.Sprintf("%s %s", quarter.Label, statusIcon)
}
