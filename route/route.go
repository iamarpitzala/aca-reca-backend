package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	_ "github.com/iamarpitzala/aca-reca-backend/docs" // swagger docs
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	payslip "github.com/iamarpitzala/aca-reca-backend/route/payship"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(e *gin.Engine) {
	cfg := config.Load()
	db, err := config.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	rdb := config.NewRedisClient(cfg.Redis)

	tokenService := service.NewTokenService(cfg.JWT)
	sessionService := service.NewSessionService(rdb, cfg.Session)
	oauthService := service.NewOAuthService(cfg.OAuth, db.DB)
	authService := service.NewAuthService(db.DB, tokenService, sessionService, oauthService)
	clinicService := service.NewClinicService(db.DB)

	authHandler := httpHandler.NewAuthHandler(authService, oauthService)
	userHandler := httpHandler.NewUserHandler(authService)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicService)

	// Swagger documentation route
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/api/v1")
	auth.RegisterAuthRoutes(v1, authHandler)
	auth.RegisterUserRoutes(v1, userHandler)
	clinic.RegisterClinicRoutes(v1, clinicHandler)
	payslip.RegisterPayslipRoutes(v1, payslipHandler)

}
