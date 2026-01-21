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

type Row struct {
	Provider    string `json:"provider"`
	Code        string `json:"code"`
	GST         string `json:"gst"`
	Collection  string `json:"collection"`
	ExternalLab string `json:"externalLab"`
	Implant     string `json:"implant"`
	InternalLab string `json:"internalLab"`
	Net         string `json:"net"`
}

type PDFConfig struct {
	ShowLogo       bool          `json:"showLogo"`
	PrimaryColor   []int         `json:"primaryColor"`   // RGB
	SecondaryColor []int         `json:"secondaryColor"` // RGB
	AccentColor    []int         `json:"accentColor"`    // RGB
	FontFamily     string        `json:"fontFamily"`
	HeaderFontSize float64       `json:"headerFontSize"`
	BodyFontSize   float64       `json:"bodyFontSize"`
	TableFontSize  float64       `json:"tableFontSize"`
	ShowSections   []string      `json:"showSections"` // ["header", "table", "reconciliation", "serviceFee", "notes"]
	CustomFields   []CustomField `json:"customFields"`
}

type PDFTableConfig struct {
}

type CustomField struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Type  string `json:"type"` // "text", "number", "currency"
}

type Statement struct {
	CompanyName            string `json:"companyName"`
	Date                   string `json:"date"`
	Reference              string `json:"reference"`
	ProviderName           string `json:"providerName"`
	SupplierCode           string `json:"supplierCode"`
	TotalPayable           string `json:"totalPayable"`
	GLCode                 string `json:"glCode"`
	Rows                   []Row  `json:"rows"`
	ServiceFee             string `json:"serviceFee"`
	TotalChargeExclGST     string `json:"totalChargeExclGST"`
	GST                    string `json:"gst"`
	TotalChargeInclGST     string `json:"totalChargeInclGST"`
	CollectedFees          string `json:"collectedFees"`
	DirectCostExternalLab  string `json:"directCostExternalLab"`
	NetCollectedFees       string `json:"netCollectedFees"`
	RemainingAmountPercent string `json:"remainingAmountPercent"`
	RAPercentNetCollected  string `json:"raPercentNetCollected"`
	GSTOnCollection        string `json:"gstOnCollection"`
	DirectCostItemsGST     string `json:"directCostItemsGST"`
	RemainingAmount        string `json:"remainingAmount"`
	GSTOnServiceFee        string `json:"gstOnServiceFee"`
	RemittedAmount         string `json:"remittedAmount"`
	ServiceFeePercent      string `json:"serviceFeePercent"`
	DentalDrawPercent      string `json:"dentalDrawPercent"`
	Notes                  string `json:"notes"`
	// Config                 *domain.PDFConfig `json:"config"`
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

func getDefaultData() Statement {
	return Statement{
		CompanyName:  "Bupa Dental",
		Date:         "08/2025",
		Reference:    "202508CHGV",
		ProviderName: "Dr Zhongjin Huang",
		SupplierCode: "ADMSC-0000003055-1",
		TotalPayable: "$16,129.78",
		GLCode:       "011000320125",
		Rows: []Row{
			{
				Provider:    "Dr Zhongjin Huang",
				Code:        "CHGV",
				GST:         "$0.00",
				Collection:  "$48,171.70",
				ExternalLab: "($665.00)",
				Implant:     "$0.00",
				InternalLab: "$0.00",
				Net:         "$47,506.70",
			},
		},
		ServiceFee:             "$28,504.02",
		TotalChargeExclGST:     "$28,504.02",
		GST:                    "$2,850.40",
		TotalChargeInclGST:     "$31,354.42",
		CollectedFees:          "$48,171.70",
		DirectCostExternalLab:  "-$665.00",
		NetCollectedFees:       "$47,506.70",
		RemainingAmountPercent: "40.0%",
		RAPercentNetCollected:  "$19,002.68",
		GSTOnCollection:        "$0.00",
		DirectCostItemsGST:     "-$22.50",
		RemainingAmount:        "$18,980.18",
		GSTOnServiceFee:        "-$2,850.40",
		RemittedAmount:         "$16,129.78",
		ServiceFeePercent:      "60.00%",
		DentalDrawPercent:      "40.00%",
		Notes:                  "",
		// Config:                 utils.GetDefaultPDFConfig(),
	}
}
