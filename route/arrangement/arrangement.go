package arrangement

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterArrangementRoutes(e *gin.RouterGroup, handler *httpHandler.ArrangmentHandler, tokenService *service.TokenService) {
	g := e.Group("/arrangements")
	g.Use(middleware.AuthMiddleware(tokenService))
	g.POST("", handler.Create)
	g.GET("", handler.ListByUserID)
	g.GET("/:id", handler.GetByID)
	g.PUT("/:id", handler.Update)
	g.DELETE("/:id", handler.Delete)
}
