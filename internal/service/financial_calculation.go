package service

import (
	"encoding/json"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type CalculationResultNet struct {
	Incomes               map[string]float64 `json:"incomes"`
	Expenses              map[string]float64 `json:"expenses"`
	NetAmount             float64            `json:"netAmount"`
	GSTCommission         float64            `json:"gstCommission"`
	CommissionForDentist  float64            `json:"commissionForDentist"`
	CommissionComponent   float64            `json:"commissionComponent"`
	SpuerComponent        float64            `json:"spuerComponent"`
	GSTOnCommission       float64            `json:"gstOnCommission"`
	TotalPayableToDentist float64            `json:"totalPayableToDentist"`
}
type GSTConfig struct {
	Rate    float64 `json:"rate"`
	Type    string  `json:"type"` // inclusive / exclusive
	Enabled bool    `json:"enabled"`
}
type CreateCustomFormRequest struct {
	ClinicID                     string          `json:"clinicId"`
	Name                         string          `json:"name"`
	Description                  string          `json:"description"`
	CalculationMethod            string          `json:"calculationMethod"`
	FormType                     string          `json:"formType"`
	Fields                       json.RawMessage `json:"fields"`
	DefaultPaymentResponsibility *string         `json:"defaultPaymentResponsibility,omitempty"`
	ServiceFacilityFeePercent    *float64        `json:"serviceFacilityFeePercent,omitempty"`
	OutworkEnabled               *bool           `json:"outworkEnabled,omitempty"`
	OutworkRatePercent           *float64        `json:"outworkRatePercent,omitempty"`
}

type Field struct {
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Label          string    `json:"label"`
	Order          int       `json:"order"`
	Section        string    `json:"section"` // income / expense
	Required       bool      `json:"required"`
	GSTConfig      GSTConfig `json:"gstConfig"`
	IncludeInTotal bool      `json:"includeInTotal"`
}

func Calculation(config domain.CreateCustomFormRequest, clinic domain.Clinic, net CalculationResultNet) string {

	switch config.CalculationMethod {
	case "gross":
		return CalculationGross(config, clinic)
	case "net":
		return CalculationNet(config, clinic, net)
	}
}

func CalculationGross(config domain.CreateCustomFormRequest, clinic domain.Clinic) string {

}

func CalculationNet(config domain.CreateCustomFormRequest, clinic domain.Clinic, net CalculationResultNet) string {
	if config.Fields == nil {
		return "No fields defined"
	}
	expen := 0.0
	income := 0.0
	net.CommissionComponent = 0.0
	net.SpuerComponent = 0.0

	for _, value := range net.Incomes {
		income += value
	}
	for _, value := range net.Expenses {
		expen += value
	}
	net.NetAmount = income - expen
	net.GSTCommission = float64(clinic.OwnerShare)

	net.CommissionForDentist = net.NetAmount * (net.GSTCommission / 100)

	if clinic.WithHoldingTax {
		net.CommissionComponent = net.CommissionForDentist / 1.12
		net.SpuerComponent = net.CommissionComponent * 0.12
		net.GSTCommission = net.CommissionComponent * 0.1
		net.TotalPayableToDentist = net.CommissionComponent + net.GSTCommission
	} else {
		net.GSTCommission = net.CommissionForDentist * 0.1
		net.TotalPayableToDentist = net.CommissionForDentist + net.GSTCommission
	}

	return "Net calculation performed for clinic: " + clinic.Name + " with net amount: " + string(net.NetAmount)
}
