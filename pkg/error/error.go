package error

var (
	ErrTaxTreatmentRequired = DomainError{
		Category: ErrorValidation,
		Code:     "TAX_TREATMENT_REQUIRED",
		Message:  "Tax treatment must be specified",
		Field:    "tax_treatment",
	}

	ErrInvalidTaxTreatment = DomainError{
		Category: ErrorValidation,
		Code:     "INVALID_TAX_TREATMENT",
		Message:  "Invalid tax treatment value",
		Field:    "tax_treatment",
	}

	ErrTaxRateRequired = DomainError{
		Category: ErrorValidation,
		Code:     "TAX_RATE_REQUIRED",
		Message:  "Tax rate is required for inclusive or exclusive tax",
		Field:    "tax_rate",
	}

	ErrTaxRateNotAllowed = DomainError{
		Category: ErrorValidation,
		Code:     "TAX_RATE_NOT_ALLOWED",
		Message:  "Tax rate must not be provided when tax treatment is none or manual",
		Field:    "tax_rate",
	}

	ErrManualTaxAmountRequired = DomainError{
		Category: ErrorValidation,
		Code:     "MANUAL_TAX_AMOUNT_REQUIRED",
		Message:  "Tax amount is required when tax treatment is manual",
		Field:    "tax_amount",
	}

	ErrManualTaxAmountNotAllowed = DomainError{
		Category: ErrorValidation,
		Code:     "MANUAL_TAX_AMOUNT_NOT_ALLOWED",
		Message:  "Tax amount must not be provided unless tax treatment is manual",
		Field:    "tax_amount",
	}

	ErrAmountsAreLocked = DomainError{
		Category: ErrorState,
		Code:     "AMOUNTS_ARE_LOCKED",
		Message:  "Cannot modify tax once journal is posted",
	}

	ErrNegativeTaxAmount = DomainError{
		Category: ErrorCalculation,
		Code:     "NEGATIVE_TAX_AMOUNT",
		Message:  "Tax amount cannot be negative",
		Field:    "tax_amount",
	}
)
