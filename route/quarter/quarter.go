package quarter

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterQuarterRoutes(e *gin.RouterGroup, userClinicHandler *httpHandler.QuarterHandler) {
	quarter := e.Group("/quarter")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	quarter.Use(middleware.AuthMiddleware(tokenService))

	quarter.POST("/", userClinicHandler.Create)
	quarter.GET("/:id", userClinicHandler.Get)
	quarter.PUT("/:id", userClinicHandler.Update)
	quarter.DELETE("/:id", userClinicHandler.Delete)
	quarter.GET("/", userClinicHandler.List)
}
