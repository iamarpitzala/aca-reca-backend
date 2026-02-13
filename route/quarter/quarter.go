package quarter

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterQuarterRoutes(e *gin.RouterGroup, quarterHandler *httpHandler.QuarterHandler, tokenService *service.TokenService) {
	// System-driven endpoints (calculate quarters from financial settings)
	clinicQuarter := e.Group("/clinic/:id")
	clinicQuarter.Use(middleware.AuthMiddleware(tokenService))

	clinicQuarter.GET("/quarters", quarterHandler.CalculateForClinic)
	clinicQuarter.GET("/quarter/date", quarterHandler.GetQuarterForDate)
}
