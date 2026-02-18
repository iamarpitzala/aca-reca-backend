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
	DefaultAmountModeGSTInclusive = "GST_INCLUSIVE"
	DefaultAmountModeGSTExclusive = "GST_EXCLUSIVE"
)

// Share / method types (clinic)
const (
	ShareTypePercentage = "PERCENTAGE"
	ShareTypeFixed      = "FIXED"
	MethodTypeNet       = "NET"
	MethodTypeGross     = "GROSS"
)

// Form type and section
const (
	FormTypeIncome   = "INCOME"
	FormTypeExpense  = "EXPENSE"
	FormTypeBoth     = "BOTH"
	SectionIncome    = "INCOME"
	SectionExpense   = "EXPENSE"
	SectionReduction = "REDUCTION"
)

// Payment responsibility (form default_payment_responsibility, field payment_responsibility)
const (
	PaymentResponsibilityOwner  = "OWNER"
	PaymentResponsibilityClinic = "CLINIC"
)

// GST type (field gst_type, tbl_gst.type)
const (
	GSTTypeInclusive = "INCLUSIVE"
	GSTTypeExclusive = "EXCLUSIVE"
	GSTTypeManual    = "MANUAL"
)

// Date formats
const (
	DateFormatRFC3339     = "2006-01-02T15:04:05Z07:00"
	DateFormatDate        = "2006-01-02"
	DateFormatYYYYMMDD   = "YYYY-MM-DD"
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
