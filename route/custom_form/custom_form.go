package custom_form

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterCustomFormRoutes(e *gin.RouterGroup, handler *httpHandler.CustomFormHandler, entryHandler *httpHandler.FieldEntryHandler, tokenService *service.TokenService) {
	g := e.Group("/custom-form")
	g.Use(middleware.AuthMiddleware(tokenService))

	g.POST("", handler.Create)
	g.GET("/clinic/:clinicId", handler.GetByClinicID)
	g.GET("/clinic/:clinicId/published", handler.GetPublishedByClinicID)
	g.GET("/:id", handler.GetByID)
	g.PUT("/:id", handler.Update)
	g.POST("/:id/publish", handler.Publish)
	g.POST("/:id/unpublish", handler.Unpublish)
	g.POST("/:id/archive", handler.Archive)
	g.DELETE("/:id", handler.Delete)
	g.POST("/:id/duplicate", handler.Duplicate)

	// Entries under /entries to avoid conflicting with form :id
	entries := g.Group("/entries")
	entries.POST("", entryHandler.Create)
	entries.GET("/form/:formId", entryHandler.GetByFormID)
	entries.GET("/clinic/:clinicId", entryHandler.GetByClinicID)
	entries.GET("/:entryId", entryHandler.GetByID)
	entries.PUT("/:entryId", entryHandler.Update)
	entries.DELETE("/:entryId", entryHandler.Delete)
}
