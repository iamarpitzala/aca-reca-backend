package error

import (
	"github.com/iamarpitzala/aca-reca-backend/util"
)

func ValidateTaxInput(treatment util.TaxTreatment, taxRate *float64, taxAmount *float64, isPosted bool) error {

	if isPosted {
		return &ErrAmountsAreLocked
	}

	switch treatment {

	case util.TaxTreatmentExclusive, util.TaxTreatmentInclusive:
		if taxRate == nil {
			return &ErrTaxRateRequired
		}
		if taxAmount != nil {
			return &ErrManualTaxAmountNotAllowed
		}

	case util.TaxTreatmentManual:
		if taxAmount == nil {
			return &ErrManualTaxAmountRequired
		}
		if taxRate != nil {
			return &ErrTaxRateNotAllowed
		}

	case util.TaxTreatmentNone:
		if taxRate != nil {
			return &ErrTaxRateNotAllowed
		}
		if taxAmount != nil {
			return &ErrManualTaxAmountNotAllowed
		}

	default:
		return &ErrInvalidTaxTreatment
	}

	return nil
}
