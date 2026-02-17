package entry

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterEntryRoutes(e *gin.RouterGroup, entryHandler *httpHandler.FieldEntryHandler, tokenService *service.TokenService) {
	entry := e.Group("/entry")
	entry.Use(middleware.AuthMiddleware(tokenService))
	entry.POST("", entryHandler.Create)
	entry.GET("/:id/net-details", entryHandler.GetNetDetails)
	entry.GET("/:id/gross-details", entryHandler.GetGrossDetails)
	entry.GET("/:id", entryHandler.GetByID)
	entry.GET("/form/:formId", entryHandler.GetByFormID)
	entry.GET("/clinic/:clinicId", entryHandler.GetByClinicID)
	entry.PUT("/:id", entryHandler.Update)
	entry.DELETE("/:id", entryHandler.Delete)
}
