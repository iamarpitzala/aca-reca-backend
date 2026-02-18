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
	"github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	clinic_financial_settings "github.com/iamarpitzala/aca-reca-backend/route/clinic_financial_settings"
	custom_form "github.com/iamarpitzala/aca-reca-backend/route/custom_form"
	"github.com/iamarpitzala/aca-reca-backend/route/entry"
	expense "github.com/iamarpitzala/aca-reca-backend/route/expense"
	payslip 	"github.com/iamarpitzala/aca-reca-backend/route/payship"
	"github.com/iamarpitzala/aca-reca-backend/route/quarter"
	bas_snapshot_route "github.com/iamarpitzala/aca-reca-backend/route/bas_snapshot"
	reports_route "github.com/iamarpitzala/aca-reca-backend/route/reports"
	transaction_route "github.com/iamarpitzala/aca-reca-backend/route/transaction"
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

	// Repositories (driven adapters)
	clinicRepo := postgres.NewClinicRepository(sqlxDB)
	clinicCOARepo := postgres.NewClinicCOARepository(sqlxDB)
	userClinicRepo := postgres.NewUserClinicRepository(sqlxDB)
	userRepo := postgres.NewUserRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	oauthProviderRepo := postgres.NewOAuthProviderRepository(sqlxDB)
	oauthService := service.NewOAuthService(cfg.OAuth, oauthProviderRepo, userRepo)
	expenseRepo := postgres.NewExpenseRepository(sqlxDB)
	aocRepo := postgres.NewAOCRepository(sqlxDB)
	customFormRepo := postgres.NewCustomFormRepository(sqlxDB)
	customFormFieldRepo := postgres.NewCustomFormFieldRepository(sqlxDB)
	customFormVersionRepo := postgres.NewCustomFormVersionRepository(sqlxDB)
	customFormCalculationRepo := postgres.NewCustomFormCalculationRepository(sqlxDB)
	clinicFinancialSettingsRepo := postgres.NewClinicFinancialSettingsRepository(sqlxDB)
	fieldEntryRepo := postgres.NewFieldEntryRepository(sqlxDB)
	entryNetDetailsRepo := postgres.NewEntryNetDetailsRepository(sqlxDB)
	entryGrossDetailsRepo := postgres.NewEntryGrossDetailsRepository(sqlxDB)
	entryGrossReductionRepo := postgres.NewEntryGrossReductionRepository(sqlxDB)
	entryGrossReimbursementRepo := postgres.NewEntryGrossReimbursementRepository(sqlxDB)
	transactionRepo := postgres.NewTransactionRepository(sqlxDB)
	pnlReportRepo := postgres.NewPnlReportRepository(sqlxDB)
	basSnapshotRepo := postgres.NewBASSnapshotRepository(sqlxDB)

	// Calculation engine (decoupled for accounting accuracy)
	calcEngine := calculation.NewEntryCalculationEngine()

	// Use cases (application layer)
	clinicUC := usecase.NewClinicService(clinicRepo)
	userClinicUC := usecase.NewUserClinicService(userClinicRepo, clinicRepo, userRepo)
	authUC := usecase.NewAuthService(userRepo, sessionRepo, tokenService, clinicUC, userClinicUC)
	clinicCOAUC := usecase.NewClinicCOAService(clinicCOARepo, clinicRepo, aocRepo)
	quarterUC := usecase.NewQuarterService(clinicFinancialSettingsRepo)
	expensesUC := usecase.NewExpensesService(expenseRepo)
	aocUC := usecase.NewAOCService(aocRepo)
	customFormUC := usecase.NewCustomFormService(customFormRepo, customFormFieldRepo, customFormVersionRepo, clinicRepo, calcEngine)
	clinicFinancialSettingsUC := usecase.NewClinicFinancialSettingsService(clinicFinancialSettingsRepo, clinicRepo)
	fieldEntryUC := usecase.NewFieldEntryService(fieldEntryRepo, customFormRepo, customFormFieldRepo, clinicRepo, customFormCalculationRepo, entryNetDetailsRepo, entryGrossDetailsRepo, entryGrossReductionRepo, entryGrossReimbursementRepo, clinicFinancialSettingsRepo, transactionRepo, calcEngine)
	transactionUC := usecase.NewTransactionService(transactionRepo, clinicRepo)
	pnlReportUC := usecase.NewPnlReportService(pnlReportRepo, clinicRepo)
	basSnapshotUC := usecase.NewBASSnapshotService(basSnapshotRepo, clinicRepo)
	// HTTP handlers (driving adapters)
	authHandler := httpHandler.NewAuthHandler(authUC, oauthService, cfg.OAuth.FrontendURL)
	userHandler := httpHandler.NewUserHandler(authUC)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicUC, userClinicUC, clinicCOAUC)
	userClinicHandler := httpHandler.NewUserClinicHandler(userClinicUC)
	customFormHandler := httpHandler.NewCustomFormHandler(customFormUC, userClinicUC)
	expensesHandler := httpHandler.NewExpensesHandler(expensesUC)
	quarterHandler := httpHandler.NewQuarterHandler(quarterUC)
	aosHandler := httpHandler.NewAOCHandler(aocUC)
	clinicFinancialSettingsHandler := httpHandler.NewClinicFinancialSettingsHandler(clinicFinancialSettingsUC)
	fieldEntryHandler := httpHandler.NewFieldEntryHandler(fieldEntryUC, userClinicUC)
	transactionHandler := httpHandler.NewTransactionHandler(transactionUC)
	pnlReportHandler := httpHandler.NewPnlReportHandler(pnlReportUC)
	basSnapshotHandler := httpHandler.NewBASSnapshotHandler(basSnapshotUC)
	// Cloudinary upload (optional: nil if env not set)
	cloudinarySvc, _ := cloudinary.NewService(cfg.Cloudinary)
	uploadHandler := httpHandler.NewUploadHandler(cloudinarySvc)

	// Swagger documentation route
	e.GET(util.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group(util.APIV1Prefix)
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
	entry.RegisterEntryRoutes(v1, fieldEntryHandler, tokenService)
	transaction_route.RegisterTransactionRoutes(v1, transactionHandler, tokenService)
	reports_route.RegisterReportsRoutes(v1, pnlReportHandler, tokenService)
	bas_snapshot_route.RegisterBASSnapshotRoutes(v1, basSnapshotHandler, tokenService)
}
