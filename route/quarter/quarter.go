package quarter

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterQuarterRoutes(e *gin.RouterGroup, quarterHandler *httpHandler.QuarterHandler, tokenService *service.TokenService) {
	quarter := e.Group("/quarter")
	quarter.Use(middleware.AuthMiddleware(tokenService))

	quarter.POST("/", quarterHandler.Create)
	quarter.GET("/:id", quarterHandler.Get)
	quarter.PUT("/:id", quarterHandler.Update)
	quarter.DELETE("/:id", quarterHandler.Delete)
	quarter.GET("/", quarterHandler.List)
}
