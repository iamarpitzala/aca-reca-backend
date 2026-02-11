package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	_ "github.com/iamarpitzala/aca-reca-backend/docs" // swagger docs
	"github.com/iamarpitzala/aca-reca-backend/internal/adapter/postgres"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/iamarpitzala/aca-reca-backend/pkg/cloudinary"
	"github.com/iamarpitzala/aca-reca-backend/route/aoc"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	custom_form "github.com/iamarpitzala/aca-reca-backend/route/custom_form"
	expense "github.com/iamarpitzala/aca-reca-backend/route/expense"
	financial_form "github.com/iamarpitzala/aca-reca-backend/route/financial_form"
	payslip "github.com/iamarpitzala/aca-reca-backend/route/payship"
	"github.com/iamarpitzala/aca-reca-backend/route/quarter"
	upload_route "github.com/iamarpitzala/aca-reca-backend/route/upload"
	user_clinic "github.com/iamarpitzala/aca-reca-backend/route/user_clinic"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(e *gin.Engine) {
	cfg := config.Load()
	db, err := config.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlxDB := db.DB

	// Token service (implements port.TokenProvider; stays in service for JWT)
	tokenService := service.NewTokenService(cfg.JWT)
	oauthService := service.NewOAuthService(cfg.OAuth, sqlxDB)

	// Repositories (driven adapters)
	clinicRepo := postgres.NewClinicRepository(sqlxDB)
	userClinicRepo := postgres.NewUserClinicRepository(sqlxDB)
	userRepo := postgres.NewUserRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	quarterRepo := postgres.NewQuarterRepository(sqlxDB)
	financialFormRepo := postgres.NewFinancialFormRepository(sqlxDB)
	expenseRepo := postgres.NewExpenseRepository(sqlxDB)
	aocRepo := postgres.NewAOCRepository(sqlxDB)

	// Use cases (application layer)
	authUC := usecase.NewAuthService(userRepo, sessionRepo, tokenService)
	clinicUC := usecase.NewClinicService(clinicRepo)
	userClinicUC := usecase.NewUserClinicService(userClinicRepo, clinicRepo, userRepo)
	quarterUC := usecase.NewQuarterService(quarterRepo)
	financialFormUC := usecase.NewFinancialFormService(financialFormRepo, clinicRepo)
	expensesUC := usecase.NewExpensesService(expenseRepo)
	aocUC := usecase.NewAOCService(aocRepo)

	// Custom form still uses legacy service (calculation + entry logic to be ported later)
	customFormService := service.NewCustomFormService(sqlxDB)

	// HTTP handlers (driving adapters)
	authHandler := httpHandler.NewAuthHandler(authUC, oauthService, cfg.OAuth.FrontendURL)
	userHandler := httpHandler.NewUserHandler(authUC)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicUC, userClinicUC)
	userClinicHandler := httpHandler.NewUserClinicHandler(userClinicUC)
	financialFormHandler := httpHandler.NewFinancialFormHandler(financialFormUC)
	customFormHandler := httpHandler.NewCustomFormHandler(customFormService)
	expensesHandler := httpHandler.NewExpensesHandler(expensesUC)
	quarterHandler := httpHandler.NewQuarterHandler(quarterUC)
	aosHandler := httpHandler.NewAOCHandler(aocUC)

	// Cloudinary upload (optional: nil if env not set)
	cloudinarySvc, _ := cloudinary.NewService(cfg.Cloudinary)
	uploadHandler := httpHandler.NewUploadHandler(cloudinarySvc)

	// Swagger documentation route
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/api/v1")
	auth.RegisterAuthRoutes(v1, authHandler)
	auth.RegisterUserRoutes(v1, userHandler, tokenService)
	clinic.RegisterClinicRoutes(v1, clinicHandler, tokenService)
	payslip.RegisterPayslipRoutes(v1, payslipHandler)
	user_clinic.RegisterUserClinicRoutes(v1, userClinicHandler, tokenService)
	financial_form.RegisterFinancialFormRoutes(v1, financialFormHandler, tokenService)
	custom_form.RegisterCustomFormRoutes(v1, customFormHandler, tokenService)
	quarter.RegisterQuarterRoutes(v1, quarterHandler, tokenService)
	expense.RegisterExpensesRoutes(v1, expensesHandler, tokenService)
	aoc.RegisterAOCRoutes(v1, aosHandler, tokenService)
	upload_route.RegisterUploadRoutes(v1, uploadHandler, tokenService)
}
