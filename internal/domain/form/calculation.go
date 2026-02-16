package form

import "github.com/iamarpitzala/aca-reca-backend/util"

type CalculationMethod string

const (
	CalculationMethodNet   CalculationMethod = CalculationMethod(util.MethodTypeNet)
	CalculationMethodGross CalculationMethod = CalculationMethod(util.MethodTypeGross)
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
