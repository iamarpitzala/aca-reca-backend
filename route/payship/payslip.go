package payslip

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
)

func RegisterPayslipRoutes(e *gin.RouterGroup, payslipHeander *httpHandler.PayslipHandler) {
	payslip := e.Group("/payslip")
	payslip.POST("/export/income", payslipHeander.ExportExcelIncome)
	payslip.POST("/export/expenses", payslipHeander.ExportExcelExpanses)
	payslip.POST("/generate/pdf", payslipHeander.GeneratePdf)
}
