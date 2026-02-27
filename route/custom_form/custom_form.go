package custom_form

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterCustomFormRoutes(e *gin.RouterGroup, handler *httpHandler.CustomFormHandler, entryHandler *httpHandler.CustomFormEntryHandler, fieldHandler *httpHandler.CustomFormFieldHandler, tokenService *service.TokenService) {
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
	// g.POST("/:id/duplicate", handler.Duplicate)

	// Form fields (by form or by version)
	forms := g.Group("/forms")
	forms.GET("/:formId/fields", fieldHandler.GetFieldsByFormID)
	forms.GET("/:formId/versions/:formVersionId/fields", fieldHandler.GetFieldsByFormVersionID)
	forms.POST("/:formId/fields", fieldHandler.CreateField)

	// Single form field (by field id)
	fields := g.Group("/fields")
	fields.GET("/:fieldId", fieldHandler.GetFieldByID)
	fields.PUT("/:fieldId", fieldHandler.UpdateField)
	fields.DELETE("/:fieldId", fieldHandler.DeleteField)

	// Field configs (nested under field)
	fields.POST("/:fieldId/configs", fieldHandler.CreateFieldConfig)
	fields.GET("/:fieldId/configs", fieldHandler.GetFieldConfigsByFormFieldID)

	// Single field config
	configs := g.Group("/field-configs")
	configs.GET("/:configId", fieldHandler.GetFieldConfigByID)
	configs.PUT("/:configId", fieldHandler.UpdateFieldConfig)
	configs.DELETE("/:configId", fieldHandler.DeleteFieldConfig)

	// Entries under /entries to avoid conflicting with form :id
	entries := g.Group("/entries")
	entries.POST("", entryHandler.Create)
	entries.GET("/form/:formId", entryHandler.GetByFormID)
	entries.GET("/clinic/:clinicId", entryHandler.GetByClinicID)
	entries.GET("/:entryId", entryHandler.GetByID)
	entries.PUT("/:entryId", entryHandler.Update)
	entries.DELETE("/:entryId", entryHandler.Delete)
}
