package util

// Role (tbl_user_clinic.role)
const (
	RoleOwner  = "OWNER"
	RoleAdmin  = "ADMIN"
	RoleMember = "MEMBER"
	RoleViewer = "VIEWER"
)

// Period / reporting frequency (clinic_financial_settings)
const (
	PeriodQuarterly = "QUARTERLY"
	PeriodAnnually  = "ANNUALLY"
)

// Financial year start
const (
	FinancialYearStartJuly    = "JULY"
	FinancialYearStartJanuary = "JANUARY"
)

// Accounting method
const (
	AccountingMethodCash    = "CASH"
	AccountingMethodAccrual = "ACCRUAL"
)

// Default amount mode
const (
	Inclusive = "INCLUSIVE"
	Exclusive = "EXCLUSIVE"
)

// Share / method types (clinic)
const (
	CalculationMethodNet   = "NET"
	CalculationMethodGross = "GROSS"
)

// Form type and section
const (
	FormTypeIncome  = "INCOME"
	FormTypeExpense = "EXPENSE"
)

// Payment responsibility (form default_payment_responsibility, field payment_responsibility)
const (
	PaymentResponsibilityOwner  = "OWNER"
	PaymentResponsibilityClinic = "CLINIC"
)

// Date formats
const (
	DateFormatRFC3339  = "2006-01-02T15:04:05Z07:00"
	DateFormatDate     = "2006-01-02"
	DateFormatYYYYMMDD = "YYYY-MM-DD"
)

// OAuth error codes (used in redirect URLs)
const (
	OAuthErrorMissingCode      = "missing_code"
	OAuthErrorExchangeFailed   = "exchange_failed"
	OAuthErrorUserInfoFailed   = "user_info_failed"
	OAuthErrorLookupFailed     = "lookup_failed"
	OAuthErrorCreateUserFailed = "create_user_failed"
	OAuthErrorLoginFailed      = "login_failed"
)

// OAuth query parameters
const (
	OAuthParamAccessToken  = "access_token"
	OAuthParamRefreshToken = "refresh_token"
	OAuthParamTokenType    = "token_type"
	OAuthParamError        = "error"
	OAuthParamCode         = "code"
)

type TaxTreatment string

const (
	TaxTreatmentInclusive TaxTreatment = "INCLUSIVE"
	TaxTreatmentExclusive TaxTreatment = "EXCLUSIVE"
	TaxTreatmentManual    TaxTreatment = "MANUAL"
	TaxTreatmentNone      TaxTreatment = "NONE"
)

func (s TaxTreatment) String() string {
	return string(s)
}

type FormStatus string

const (
	FormStatusDraft     FormStatus = "DRAFT"
	FormStatusPublished FormStatus = "PUBLISHED"
	FormStatusArchived  FormStatus = "ARCHIVED"
)

func (s FormStatus) String() string {
	return string(s)
}

type FormSection string

const (
	FormSectionIncome  FormSection = "INCOME"
	FormSectionExpense FormSection = "EXPENSE"
)

func (s FormSection) String() string {
	return string(s)
}

type AccountingMethod string

const (
	AccountingMethodNet   AccountingMethod = "NET"
	AccountingMethodGross AccountingMethod = "GROSS"
)

func (s AccountingMethod) String() string {
	return string(s)
}

type ShareType string

const (
	ShareTypePercentage ShareType = "PERCENTAGE"
	ShareTypeFixed      ShareType = "FIXED"
)

func (s ShareType) String() string {
	return string(s)
}
