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

// gross calculation

// 1st step  : calculate  patient_fee(excl GST)
// in label to get lable id and label id define  suppose  GST true always manul GST entry
// if(type == income)
// if(GST) patient_fee(excl GST) = key.amount - key.gstamount      else atient_fee(excl GST) = key.amount

// 2nd step calculate expenses
// if(type == expenses)
// in each expenses
// chack label id to what gst  (inclusive , exculsive , manul , gst free)  and  payby(cinic or denties)
// suppose  inclusive amount then     net expenses will be  expenses.amount - amount/11  and gst += amount/11
// and GST add on in remmitted section
// in this section two types of expenses tracking  total_net_expenses  and  pay_by_denties

// 3rd step calculate net amount
//  net_amoiunt = patient_fee(excl GST) - total_net_expenses

// 4th step calculate total service and facility fee with other cost if on if outword cost true
//if (outword_cost) service_and_facility_fee = patient_fee(excl GST)* 60 ;
// lab_fee_wuth_OtherCostCharge  =  (total_net_expenses + GST) * 40
// total_service_andFaclilty_with_othercost = lab_fee_wuth_OtherCostCharge +service_and_facility_fee
// GST_on_service_feeFee = total_service_andFaclilty_with_othercost*10
//total_service_andFaclilty_with_othercost_with_inc_GST = total_service_andFaclilty_with_othercost + GST_on_service_feeFee

// else  service_and_facility_fee =  net_amount * 60
// gst_on_service_and_facility_fee  = service_and_facility_fee*10
//total_service_and_facility_fee = service_and_facility_fee + gst_on_service_and_facility_fee

// step 5th calculate other rimmited cost
//(type == remitted and !outword_cost) // chack label id to what gst  (inclusive , exculsive , manul , gst free)
// suppose  inclusive amount then     net expenses will be  expenses.amount - amount/11  and gst += amount/11
//  inclusive , total remitted = expenses.amount - amount/11
// exclusive , total remitted = expenses.amount  + expenses.amount * 0.1
// alos 2nd step expenses GST also show inside this only if on if pay by clinic

// Add editional reambussion
// pay by denties all expenses with gst show

// step 6th total_amount_payavle_to_denties
//if(outword) net_amount -( total_service_and_faclilty_fee_with_outwordcost)
//

//  net_amount -( total_service_and_faclilty_fee  + remittedCost) + if(paybyDenties).(net_Expenses + GST_on_this_expenses)

func CalculationGross(customForm domain.CustomForm, clinic *domain.Clinic) {
}

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
