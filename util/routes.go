package util

// API route paths
const (
	APIV1Prefix = "/api/v1"
	SwaggerPath = "/swagger/*any"
)

// Route group paths
const (
	RouteAuth                    = "/auth"
	RouteUser                    = "/user"
	RouteClinic                  = "/clinic"
	RoutePayslip                 = "/payslip"
	RouteUserClinic              = "/user-clinic"
	RouteCustomForm              = "/custom-form"
	RouteQuarter                 = "/quarter"
	RouteExpense                 = "/expense"
	RouteAOC                     = "/coa"
	RouteUpload                  = "/upload"
	RouteClinicFinancialSettings = "/clinic-financial-settings"
	RouteEntry                   = "/entry"
	RouteTransaction             = "/transaction"
	RouteReports                 = "/reports"
)

// Upload route paths
const (
	RouteUploadImage    = "/upload/image"
	RouteUploadDocument = "/upload/document"
)
