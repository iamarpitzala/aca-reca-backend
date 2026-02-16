package util

// Role (tbl_user_clinic.role)
const (
	RoleOwner  = "OWNER"
	RoleAdmin  = "ADMIN"
	RoleMember = "MEMBER"
	RoleViewer = "VIEWER"
)

// Period / reporting frequency (clinic_financial_settings, bas_snapshot)
const (
	PeriodQuarterly = "QUARTERLY"
	PeriodAnnually  = "ANNUALLY"
)

// BAS snapshot status
const (
	BASStatusDraft     = "DRAFT"
	BASStatusFinalised = "FINALISED"
	BASStatusLocked    = "LOCKED"
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
