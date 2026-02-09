package custom_form

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterCustomFormRoutes(e *gin.RouterGroup, handler *httpHandler.CustomFormHandler) {
	g := e.Group("/custom-form")
	cfg := config.Load()
	tokenService := service.NewTokenService(cfg.JWT)
	g.Use(middleware.AuthMiddleware(tokenService))

	g.POST("", handler.Create)
	g.GET("/clinic/:clinicId", handler.GetByClinicID)
	g.GET("/clinic/:clinicId/published", handler.GetPublishedByClinicID)
	g.GET("/:id", handler.GetByID)
	g.PUT("/:id", handler.Update)
	g.POST("/:id/publish", handler.Publish)
	g.POST("/:id/archive", handler.Archive)
	g.DELETE("/:id", handler.Delete)
	g.POST("/:id/duplicate", handler.Duplicate)

	// Entries under /entries to avoid conflicting with form :id
	entries := g.Group("/entries")
	entries.POST("/preview", handler.PreviewCalculations)
	entries.POST("", handler.CreateEntry)
	entries.GET("/form/:formId", handler.GetEntriesByFormID)
	entries.GET("/clinic/:clinicId", handler.GetEntriesByClinicID)
	entries.GET("/:entryId", handler.GetEntryByID)
	entries.PUT("/:entryId", handler.UpdateEntry)
	entries.DELETE("/:entryId", handler.DeleteEntry)
}
