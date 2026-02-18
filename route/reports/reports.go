package reports

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterReportsRoutes(e *gin.RouterGroup, pnlHandler *httpHandler.PnlReportHandler, tokenService *service.TokenService) {
	reports := e.Group("/reports")
	reports.Use(middleware.AuthMiddleware(tokenService))

	pnl := reports.Group("/pnl")
	pnl.POST("/generate", pnlHandler.Generate)
	pnl.GET("/:id", pnlHandler.GetByID)
	pnl.GET("/clinic/:clinicId", pnlHandler.GetByClinicID)
	pnl.POST("/:id/finalize", pnlHandler.Finalize)
	pnl.POST("/:id/regenerate", pnlHandler.Regenerate)
	pnl.DELETE("/:id", pnlHandler.Delete)
}
