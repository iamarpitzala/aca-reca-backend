package financial_calculation

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterFinancialCalculationRoutes(e *gin.RouterGroup, calculationHandler *httpHandler.FinancialCalculationHandler) {
	calculation := e.Group("/financial-calculation")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	calculation.Use(middleware.AuthMiddleware(tokenService))

	calculation.POST("/calculate", calculationHandler.CalculateFinancial)
	calculation.GET("/history/:formId", calculationHandler.GetCalculationHistory)
}
