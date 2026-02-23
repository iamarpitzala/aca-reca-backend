package error

type ErrorCategory string

const (
	ErrorValidation  ErrorCategory = "VALIDATION"
	ErrorCalculation ErrorCategory = "CALCULATION"
	ErrorState       ErrorCategory = "STATE"
)

type DomainError struct {
	Category ErrorCategory  `json:"category"`
	Code     string         `json:"code"`
	Message  string         `json:"message"`
	Field    string         `json:"field,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

func (e *DomainError) Error() string {
	return e.Code + ": " + e.Message
}
