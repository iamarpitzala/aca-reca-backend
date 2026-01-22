package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// ValidateState validates that the state is a valid Australian state/territory
func ValidateState(state string) error {
	state = strings.ToUpper(strings.TrimSpace(state))
	for _, validState := range domain.ValidStates {
		if state == validState {
			return nil
		}
	}
	return errors.New("invalid state. Must be one of: NSW, VIC, QLD, SA, WA, TAS, NT, ACT")
}

// ValidateABN validates ABN format (11 digits)
func ValidateABN(abn string) error {
	// Remove spaces and hyphens
	abn = strings.ReplaceAll(abn, " ", "")
	abn = strings.ReplaceAll(abn, "-", "")

	// Check if it's exactly 11 digits
	matched, err := regexp.MatchString(`^\d{11}$`, abn)
	if err != nil {
		return errors.New("error validating ABN format")
	}

	if !matched {
		return errors.New("ABN must be exactly 11 digits")
	}

	return nil
}

// ValidateCalculationInput validates input based on form configuration
func ValidateCalculationInput(config map[string]interface{}, input domain.CalculationInput, method string) error {
	if input.GrossPatientFee < 0 {
		return errors.New("gross patient fee cannot be negative")
	}

	if method == "net" {
		return validateNetMethodInput(config, input)
	} else if method == "gross" {
		return validateGrossMethodInput(config, input)
	}

	return errors.New("invalid calculation method")
}

func validateNetMethodInput(config map[string]interface{}, input domain.CalculationInput) error {
	// Check if lab fee is enabled
	if labFeeEnabled, ok := config["labFeeEnabled"].(bool); ok && labFeeEnabled {
		if input.LabFee < 0 {
			return errors.New("lab fee cannot be negative")
		}
	}

	return nil
}

func validateGrossMethodInput(config map[string]interface{}, input domain.CalculationInput) error {
	// B1: Standard - check lab fee if enabled
	if labFeeEnabled, ok := config["labFeeEnabled"].(bool); ok && labFeeEnabled {
		if labFeePaidBy, ok := config["labFeePaidBy"].(string); ok && labFeePaidBy == "clinic" {
			if input.LabFee < 0 {
				return errors.New("lab fee cannot be negative")
			}
		}
	}

	// B2: With GST on Lab Fee
	if gstOnLabFee, ok := config["gstOnLabFee"].(bool); ok && gstOnLabFee {
		if input.GSTOnLabFee < 0 {
			return errors.New("GST on lab fee cannot be negative")
		}
	}

	// B3: With Merchant Fee
	if merchantFeeEnabled, ok := config["merchantFeeEnabled"].(bool); ok && merchantFeeEnabled {
		if input.MerchantFeeInclGST < 0 {
			return errors.New("merchant fee cannot be negative")
		}
		if input.BankFee < 0 {
			return errors.New("bank fee cannot be negative")
		}
	}

	// B4: GST on Patient Fee + Lab Fee Paid by Dentist
	if gstOnPatientFee, ok := config["gstOnPatientFee"].(bool); ok && gstOnPatientFee {
		if labFeePaidBy, ok := config["labFeePaidBy"].(string); ok && labFeePaidBy == "dentist" {
			if input.LabFeePaidByDentist < 0 {
				return errors.New("lab fee paid by dentist cannot be negative")
			}
		}
	}

	// B5: Outwork Charge Rate
	if outworkChargeRateEnabled, ok := config["outworkChargeRateEnabled"].(bool); ok && outworkChargeRateEnabled {
		if input.OutworkLabFee < 0 {
			return errors.New("outwork lab fee cannot be negative")
		}
		if input.OutworkMerchantFee < 0 {
			return errors.New("outwork merchant fee cannot be negative")
		}
		if input.OutworkGSTOnLabFee < 0 {
			return errors.New("outwork GST on lab fee cannot be negative")
		}
		if input.OutworkGSTOnMerchantFee < 0 {
			return errors.New("outwork GST on merchant fee cannot be negative")
		}
	}

	return nil
}
