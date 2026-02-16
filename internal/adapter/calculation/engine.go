package calculation

import (
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/calculation"
)

// EngineAdapter implements port.EntryCalculationEngine using the internal calculation package.
type EngineAdapter struct{}

func NewEntryCalculationEngine() port.EntryCalculationEngine {
	return &EngineAdapter{}
}

func (e *EngineAdapter) RunEntryCalculation(
	formFieldsJSON []byte,
	formType string,
	formCalculationMethod string,
	serviceFacilityFeePercent *float64,
	outworkEnabled bool,
	outworkRatePercent *float64,
	valuesJSON, deductionsJSON []byte,
) ([]byte, error) {
	return calculation.RunEntryCalculation(
		formFieldsJSON,
		formType,
		formCalculationMethod,
		serviceFacilityFeePercent,
		outworkEnabled,
		outworkRatePercent,
		valuesJSON,
		deductionsJSON,
	)
}
