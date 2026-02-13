package bas_snapshot

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterBASSnapshotRoutes(e *gin.RouterGroup, handler *httpHandler.BASSnapshotHandler, tokenService *service.TokenService) {
	// Clinic-specific routes (more specific routes first)
	clinicBAS := e.Group("/clinic/:id/bas-snapshots")
	clinicBAS.Use(middleware.AuthMiddleware(tokenService))
	clinicBAS.GET("", handler.GetBASSnapshotsByClinic)
	
	clinicBASCreate := e.Group("/clinic/:id/bas-snapshot")
	clinicBASCreate.Use(middleware.AuthMiddleware(tokenService))
	clinicBASCreate.POST("", handler.CreateBASSnapshot)
	
	// BAS snapshot routes
	bas := e.Group("/bas-snapshot")
	bas.Use(middleware.AuthMiddleware(tokenService))
	
	bas.GET("/:id", handler.GetBASSnapshot)
	bas.PUT("/:id", handler.UpdateBASSnapshot)
	bas.POST("/:id/finalise", handler.FinaliseBAS)
	bas.POST("/:id/lock", handler.LockBAS)
	
	// Consolidated GST Summary (management view)
	reports := e.Group("/reports")
	reports.Use(middleware.AuthMiddleware(tokenService))
	reports.POST("/consolidated-gst-summary", handler.GetConsolidatedGSTSummary)
}
