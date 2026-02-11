package form

type CalculationMethod string

const (
	CalculationMethodNet   CalculationMethod = "NET"
	CalculationMethodGross CalculationMethod = "GROSS"
)

func (c CalculationMethod) String() string {
	return string(c)
}

func (c CalculationMethod) ToCalculationMethod() CalculationMethod {
	switch c {
	case CalculationMethodNet:
		return CalculationMethodNet
	case CalculationMethodGross:
		return CalculationMethodGross
	default:
		return CalculationMethodNet
	}
}
