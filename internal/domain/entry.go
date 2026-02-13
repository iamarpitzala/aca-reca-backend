package domain

type CommonEntry struct {
	Incomes  map[string]float64 `json:"incomes"`
	Expenses map[string]float64 `json:"expenses"`
}

type CalculationResultNet struct {
	Fields                []CommonEntry `json:"commonEntry"`
	NetAmount             float64       `json:"netAmount"`
	GSTCommission         float64       `json:"gstCommission"`
	CommissionForDentist  float64       `json:"commissionForDentist"`
	CommissionComponent   float64       `json:"commissionComponent"`
	SpuerComponent        float64       `json:"spuerComponent"`
	GSTOnCommission       float64       `json:"gstOnCommission"`
	TotalPayableToDentist float64       `json:"totalPayableToDentist"`
}
