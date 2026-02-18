package bas_snapshot

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterBASSnapshotRoutes(e *gin.RouterGroup, basHandler *httpHandler.BASSnapshotHandler, tokenService *service.TokenService) {
	// Clinic-scoped routes: use full path with :id to match existing clinic wildcard
	clinicBAS := e.Group("/clinic/:id/bas-snapshot")
	clinicBAS.Use(middleware.AuthMiddleware(tokenService))
	clinicBAS.POST("", basHandler.Create)
	clinicBAS.POST("/generate", basHandler.Generate)

	clinicBASList := e.Group("/clinic/:id/bas-snapshots")
	clinicBASList.Use(middleware.AuthMiddleware(tokenService))
	clinicBASList.GET("", basHandler.GetByClinicID)

	// Snapshot-scoped routes
	bas := e.Group("/bas-snapshot")
	bas.Use(middleware.AuthMiddleware(tokenService))
	bas.GET("/:id", basHandler.GetByID)
	bas.PUT("/:id", basHandler.Update)
	bas.POST("/:id/finalise", basHandler.Finalise)
	bas.POST("/:id/lock", basHandler.Lock)
	bas.DELETE("/:id", basHandler.Delete)
}
