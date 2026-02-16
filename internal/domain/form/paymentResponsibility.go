package form

import "github.com/iamarpitzala/aca-reca-backend/util"

type PaymentResponsibility string

const (
	PaymentResponsibilityOwner  PaymentResponsibility = PaymentResponsibility(util.PaymentResponsibilityOwner)
	PaymentResponsibilityClinic PaymentResponsibility = PaymentResponsibility(util.PaymentResponsibilityClinic)
)

func (p PaymentResponsibility) String() string {
	return string(p)
}
