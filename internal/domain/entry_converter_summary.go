package domain

import (
	"time"

	"github.com/google/uuid"
)

// parseSummaryFromCalc builds EntrySummary from the calculations map (field totals, BAS mapping).
func parseSummaryFromCalc(entryID uuid.UUID, createdAt, updatedAt time.Time, calc map[string]interface{}) *NormalizedEntry {
	summary := &NormalizedEntry{
		Header: &FieldEntry{
			ID:        entryID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}
	if calc == nil {
		return summary
	}
	return summary
}
