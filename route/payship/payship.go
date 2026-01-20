package payslip

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
)

func RegisterPayshipRoutes(e *gin.RouterGroup, payslipHeander *httpHandler.PayslipHandler) {
	payslip := e.Group("/payslip/export")
	payslip.POST("/income", payslipHeander.ExportExcelIncome)
	payslip.POST("/expanses", payslipHeander.ExportExcelExpanses)
}
