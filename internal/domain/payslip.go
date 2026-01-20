package domain

type ExportIncome struct {
	Q                 string  `json:"q"`
	IncomeType        string  `json:"income_type"`
	LabRecord         string  `json:"lab_record"`
	PaymentDate       string  `json:"payment_date"`
	DentalPractice    string  `json:"dental_practice"`
	Adjustments       float64 `json:"adjustments"`
	GrossIncomeG1     float64 `json:"gross_income_g1"`
	LabFees           float64 `json:"lab_fees"`
	GrossNetLabFees   float64 `json:"gross_net_lab_fees"`
	GSTPayable1A      float64 `json:"gst_payable_1a"`
	GSTFree           float64 `json:"gst_free"`
	ManagementFeesG11 float64 `json:"management_fees_g11"`
	Percentage        float64 `json:"percentage"`
	GSTRefundable1B   float64 `json:"gst_refundable_1b"`
	NetPayment        float64 `json:"net_payment"`
}

type ExportExpenses struct {
	Q        string  `json:"q"`
	Date     string  `json:"date"`
	Supplier string  `json:"supplier"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Bas      float64 `json:"bas"`
	Net      float64 `json:"net"`
	GST      float64 `json:"gst"`
	Remarks  string  `json:"remarks"`
}

var HeadersExpenses = []string{
	"Q",
	"Payment Date",
	"Supplier Name",
	"Category*",
	"Amount",
	"% Bus use",
	"Net Amount",
	"GST",
	"Remarks",
}

var HeadersIncome = []string{
	"Q",
	"Income Type",
	"Lab Record",
	"Payment Date",
	"Dental Practice",
	"Adjustments",
	"(G1) Gross Income",
	"Lab Fees",
	"Gross Income net of Lab fees",
	"(1A) GST Payable",
	"(GST Free)",
	"(G11) Management fees",
	"%",
	"(1B) GST Refundable",
	"Net Payment",
}
