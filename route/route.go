package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	_ "github.com/iamarpitzala/aca-reca-backend/docs" // swagger docs
	"github.com/iamarpitzala/aca-reca-backend/internal/adapter/calculation"
	"github.com/iamarpitzala/aca-reca-backend/internal/adapter/postgres"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/iamarpitzala/aca-reca-backend/pkg/cloudinary"
	"github.com/iamarpitzala/aca-reca-backend/route/aoc"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	bas_snapshot "github.com/iamarpitzala/aca-reca-backend/route/bas_snapshot"
	clinic_financial_settings "github.com/iamarpitzala/aca-reca-backend/route/clinic_financial_settings"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	custom_form "github.com/iamarpitzala/aca-reca-backend/route/custom_form"
	expense "github.com/iamarpitzala/aca-reca-backend/route/expense"
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
	clinicCOARepo := postgres.NewClinicCOARepository(sqlxDB)
	userClinicRepo := postgres.NewUserClinicRepository(sqlxDB)
	userRepo := postgres.NewUserRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	expenseRepo := postgres.NewExpenseRepository(sqlxDB)
	aocRepo := postgres.NewAOCRepository(sqlxDB)
	customFormRepo := postgres.NewCustomFormRepository(sqlxDB)
	transactionRepo := postgres.NewTransactionRepository(sqlxDB)
		clinicFinancialSettingsRepo := postgres.NewClinicFinancialSettingsRepository(sqlxDB)
		basSnapshotRepo := postgres.NewBASSnapshotRepository(sqlxDB)

	// Calculation engine (decoupled for accounting accuracy)
	calcEngine := calculation.NewEntryCalculationEngine()

	// Use cases (application layer)
	authUC := usecase.NewAuthService(userRepo, sessionRepo, tokenService)
	clinicUC := usecase.NewClinicService(clinicRepo)
	clinicCOAUC := usecase.NewClinicCOAService(clinicCOARepo, clinicRepo, aocRepo)
	userClinicUC := usecase.NewUserClinicService(userClinicRepo, clinicRepo, userRepo)
	quarterUC := usecase.NewQuarterService(clinicFinancialSettingsRepo)
	expensesUC := usecase.NewExpensesService(expenseRepo)
	aocUC := usecase.NewAOCService(aocRepo)
		customFormUC := usecase.NewCustomFormService(customFormRepo, clinicRepo, calcEngine)
		transactionPostingUC := usecase.NewTransactionPostingService(customFormRepo, transactionRepo, clinicCOARepo, aocRepo)
		clinicFinancialSettingsUC := usecase.NewClinicFinancialSettingsService(clinicFinancialSettingsRepo, clinicRepo)
		basSnapshotUC := usecase.NewBASSnapshotService(basSnapshotRepo, clinicRepo)

	// HTTP handlers (driving adapters)
	authHandler := httpHandler.NewAuthHandler(authUC, oauthService, cfg.OAuth.FrontendURL)
	userHandler := httpHandler.NewUserHandler(authUC)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicUC, userClinicUC, clinicCOAUC)
	userClinicHandler := httpHandler.NewUserClinicHandler(userClinicUC)
	customFormHandler := httpHandler.NewCustomFormHandler(customFormUC, transactionPostingUC, userClinicUC)
	expensesHandler := httpHandler.NewExpensesHandler(expensesUC)
	quarterHandler := httpHandler.NewQuarterHandler(quarterUC)
	aosHandler := httpHandler.NewAOCHandler(aocUC)
		clinicFinancialSettingsHandler := httpHandler.NewClinicFinancialSettingsHandler(clinicFinancialSettingsUC)
		basSnapshotHandler := httpHandler.NewBASSnapshotHandler(basSnapshotUC)

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
	custom_form.RegisterCustomFormRoutes(v1, customFormHandler, tokenService)
	quarter.RegisterQuarterRoutes(v1, quarterHandler, tokenService)
	expense.RegisterExpensesRoutes(v1, expensesHandler, tokenService)
	aoc.RegisterAOCRoutes(v1, aosHandler, tokenService)
	upload_route.RegisterUploadRoutes(v1, uploadHandler, tokenService)
		clinic_financial_settings.RegisterClinicFinancialSettingsRoutes(v1, clinicFinancialSettingsHandler, tokenService)
		bas_snapshot.RegisterBASSnapshotRoutes(v1, basSnapshotHandler, tokenService)
}
