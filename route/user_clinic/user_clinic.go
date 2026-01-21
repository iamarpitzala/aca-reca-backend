package user_clinic

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterUserClinicRoutes(e *gin.RouterGroup, userClinicHandler *httpHandler.UserClinicHandler) {
	userClinic := e.Group("/user-clinic")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	userClinic.Use(middleware.AuthMiddleware(tokenService))

	userClinic.POST("/", userClinicHandler.AssociateUserWithClinic)
	userClinic.GET("/user/:userId", userClinicHandler.GetUserClinics)
	userClinic.GET("/clinic/:clinicId", userClinicHandler.GetClinicUsers)
	userClinic.DELETE("/:id", userClinicHandler.RemoveUserFromClinic)
}
