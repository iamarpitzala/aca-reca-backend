package clinic_financial_settings

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterClinicFinancialSettingsRoutes(e *gin.RouterGroup, handler *httpHandler.ClinicFinancialSettingsHandler, tokenService *service.TokenService) {
	settings := e.Group("/clinic/:id/financial-settings")
	settings.Use(middleware.AuthMiddleware(tokenService))
	
	settings.GET("", handler.GetFinancialSettings)
	settings.PUT("", handler.CreateOrUpdateFinancialSettings)
}
