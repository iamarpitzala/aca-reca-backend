package service

import (
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// type GSTConfig struct {
// 	Rate    float64 `json:"rate"`
// 	Type    string  `json:"type"` // inclusive / exclusive
// 	Enabled bool    `json:"enabled"`
// }
// type CreateCustomFormRequest struct {
// 	ClinicID                     string          `json:"clinicId"`
// 	Name                         string          `json:"name"`
// 	Description                  string          `json:"description"`
// 	CalculationMethod            string          `json:"calculationMethod"`
// 	FormType                     string          `json:"formType"`
// 	Fields                       json.RawMessage `json:"fields"`
// 	DefaultPaymentResponsibility *string         `json:"defaultPaymentResponsibility,omitempty"`
// 	ServiceFacilityFeePercent    *float64        `json:"serviceFacilityFeePercent,omitempty"`
// 	OutworkEnabled               *bool           `json:"outworkEnabled,omitempty"`
// 	OutworkRatePercent           *float64        `json:"outworkRatePercent,omitempty"`
// }

// type Field struct {
// 	Name           string    `json:"name"`
// 	Type           string    `json:"type"`
// 	Label          string    `json:"label"`
// 	Order          int       `json:"order"`
// 	Section        string    `json:"section"` // income / expense
// 	Required       bool      `json:"required"`
// 	GSTConfig      GSTConfig `json:"gstConfig"`
// 	IncludeInTotal bool      `json:"includeInTotal"`
// }

type CalculateGross struct {
	Fields []domain.CommonEntry `json:"commonEntry"`
}

func CommonCalculation(
	customForm *domain.CustomForm,
	clinic *domain.Clinic,
	commonEntry domain.CommonEntry,
) *domain.CalculationResultNet {

	switch customForm.CalculationMethod {
	// case "gross":
	// 	return CalculationGross(customForm, clinic, commonEntry)
	case "net":
		return CalculationNet(customForm, clinic, commonEntry)
	default:
		return nil
	}
}

// func CalculationGross(customForm domain.CustomForm, clinic *domain.Clinic) {
// 	return nil
// }

func CalculationNet(customForm *domain.CustomForm, clinic *domain.Clinic, commonEntry domain.CommonEntry) *domain.CalculationResultNet {

	var Calculation domain.CalculationResultNet
	if customForm.Fields == nil {
		return nil
	}
	expenses := 0.0
	income := 0.0
	Calculation.CommissionComponent = 0.0
	Calculation.SuperComponent = 0.0

	for _, value := range commonEntry.Incomes {
		income += value
	}
	for _, value := range commonEntry.Expenses {
		expenses += value
	}
	Calculation.NetAmount = income - expenses
	Calculation.GSTCommission = float64(clinic.OwnerShare)

	Calculation.CommissionForDentist = Calculation.NetAmount * (Calculation.GSTCommission / 100)

	if clinic.WithHoldingTax {
		Calculation.CommissionComponent = Calculation.CommissionForDentist / 1.12
		Calculation.SuperComponent = Calculation.CommissionComponent * 0.12
		Calculation.GSTOnCommission = Calculation.CommissionComponent * 0.1
		Calculation.TotalPaybleToDentist = Calculation.CommissionComponent + Calculation.GSTCommission
	} else {
		Calculation.GSTOnCommission = Calculation.CommissionForDentist * 0.1
		Calculation.TotalPaybleToDentist = Calculation.CommissionForDentist + Calculation.GSTCommission
	}

	return &Calculation
}
