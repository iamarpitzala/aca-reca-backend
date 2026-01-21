package clinic

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterClinicRoutes(e *gin.RouterGroup, clinicHandler *httpHandler.ClinicHandler) {
	clinic := e.Group("/clinic")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	clinic.Use(middleware.AuthMiddleware(tokenService))

	clinic.POST("/", clinicHandler.CreateClinic)
	clinic.GET("/:id", clinicHandler.GetClinic)
	clinic.PUT("/:id", clinicHandler.UpdateClinic)
	clinic.DELETE("/:id", clinicHandler.DeleteClinic)
	clinic.GET("/", clinicHandler.GetAllClinics)
	clinic.GET("/abn/:abnNumber", clinicHandler.GetClinicByABNNumber)
}
