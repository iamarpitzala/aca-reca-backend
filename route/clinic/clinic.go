package clinic

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterClinicRoutes(e *gin.RouterGroup, clinicHandler *httpHandler.ClinicHandler, tokenService *service.TokenService) {
	clinic := e.Group("/clinic")
	clinic.Use(middleware.AuthMiddleware(tokenService))

	clinic.POST("", clinicHandler.CreateClinic)
	clinic.GET("", clinicHandler.GetAllClinics)
	clinic.GET("/abn/:abnNumber", clinicHandler.GetClinicByABNNumber)
	// More specific paths first so /:id does not capture "id/activate" or "id/deactivate"
	clinic.PATCH("/:id/activate", clinicHandler.ActivateClinic)
	clinic.PATCH("/:id/deactivate", clinicHandler.DeactivateClinic)
	clinic.GET("/:id", clinicHandler.GetClinic)
	clinic.PUT("/:id", clinicHandler.UpdateClinic)
	clinic.DELETE("/:id", clinicHandler.DeleteClinic)
}
