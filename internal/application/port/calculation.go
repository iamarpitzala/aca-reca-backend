package port

// EntryCalculationEngine computes entry totals (base, GST, net) from form values and deductions.
// Decouples calculation logic from use case layer for testability and accounting accuracy.
type EntryCalculationEngine interface {
	RunEntryCalculation(
		formFieldsJSON []byte,
		formType string,
		serviceFacilityFeePercent *float64,
		outworkEnabled bool,
		outworkRatePercent *float64,
		valuesJSON, deductionsJSON []byte,
	) ([]byte, error)
}
