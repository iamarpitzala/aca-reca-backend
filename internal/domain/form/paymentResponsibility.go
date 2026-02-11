package form

type PaymentResponsibility string

const (
	PaymentResponsibilityOwner  PaymentResponsibility = "OWNER"
	PaymentResponsibilityClinic PaymentResponsibility = "CLINIC"
)

func (p PaymentResponsibility) String() string {
	return string(p)
}
