package user_clinic

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterUserClinicRoutes(e *gin.RouterGroup, userClinicHandler *httpHandler.UserClinicHandler, tokenService *service.TokenService) {
	userClinic := e.Group("/user-clinic")
	userClinic.Use(middleware.AuthMiddleware(tokenService))

	userClinic.POST("/", userClinicHandler.AssociateUserWithClinic)
	userClinic.GET("/user/:userId", userClinicHandler.GetUserClinics)
	userClinic.GET("/clinic/:clinicId", userClinicHandler.GetClinicUsers)
	userClinic.DELETE("/:id", userClinicHandler.RemoveUserFromClinic)
}
