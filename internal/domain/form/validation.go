package form

type ValidationRule string

const (
	ValidationRuleRequired ValidationRule = "REQUIRED"
	ValidationRuleMin      ValidationRule = "MIN"
	ValidationRuleMax      ValidationRule = "MAX"
	ValidationRulePattern  ValidationRule = "PATTERN"
	ValidationRuleCustom   ValidationRule = "CUSTOM"
)

func (v ValidationRule) String() string {
	return string(v)
}
