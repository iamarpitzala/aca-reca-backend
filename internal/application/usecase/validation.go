package usecase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// ValidateState validates that the state is a valid Australian state/territory.
func ValidateState(state string) error {
	state = strings.ToUpper(strings.TrimSpace(state))
	for _, validState := range domain.ValidStates {
		if state == validState {
			return nil
		}
	}
	return errors.New("invalid state. Must be one of: NSW, VIC, QLD, SA, WA, TAS, NT, ACT")
}

// ValidateABN validates ABN format (11 digits).
func ValidateABN(abn string) error {
	abn = strings.ReplaceAll(abn, " ", "")
	abn = strings.ReplaceAll(abn, "-", "")
	matched, err := regexp.MatchString(`^\d{11}$`, abn)
	if err != nil {
		return errors.New("error validating ABN format")
	}
	if !matched {
		return errors.New("ABN must be exactly 11 digits")
	}
	return nil
}
