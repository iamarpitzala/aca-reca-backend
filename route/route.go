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
	"github.com/iamarpitzala/aca-reca-backend/route/arrangement"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	clinic_financial_settings "github.com/iamarpitzala/aca-reca-backend/route/clinic_financial_settings"
	coa_route "github.com/iamarpitzala/aca-reca-backend/route/coa"
	custom_form "github.com/iamarpitzala/aca-reca-backend/route/custom_form"
	"github.com/iamarpitzala/aca-reca-backend/route/entry"
	expense "github.com/iamarpitzala/aca-reca-backend/route/expense"
	payslip "github.com/iamarpitzala/aca-reca-backend/route/payslip"
	reports_route "github.com/iamarpitzala/aca-reca-backend/route/reports"
	transaction_route "github.com/iamarpitzala/aca-reca-backend/route/transaction"
	upload_route "github.com/iamarpitzala/aca-reca-backend/route/upload"
	user_clinic "github.com/iamarpitzala/aca-reca-backend/route/user_clinic"
	"github.com/iamarpitzala/aca-reca-backend/util"
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
	userClinicRepo := postgres.NewUserClinicRepository(sqlxDB)
	userRepo := postgres.NewUserRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	oauthProviderRepo := postgres.NewOAuthProviderRepository(sqlxDB)
	oauthService := service.NewOAuthService(cfg.OAuth, oauthProviderRepo, userRepo)
	expenseRepo := postgres.NewExpenseRepository(sqlxDB)
	customFormRepo := postgres.NewCustomFormRepository(sqlxDB)
	customFormFieldRepo := postgres.NewCustomFormFieldRepository(sqlxDB)
	customFormVersionRepo := postgres.NewCustomFormVersionRepository(sqlxDB)
	financialYearRepo := postgres.NewFinancialYearRepository(sqlxDB)
	financialQuarterRepo := postgres.NewFinancialQuarterRepository(sqlxDB)
	clinicFinancialYearRepo := postgres.NewClinicFinancialYearRepository(sqlxDB)
	clinicFinancialQuarterLockRepo := postgres.NewClinicFinancialQuarterLockRepository(sqlxDB)
	clinicFinancialSettingsRepo := postgres.NewClinicFinancialSettingRepository(sqlxDB)

	fieldEntryRepo := postgres.NewFieldEntryRepository(sqlxDB)

	transactionRepo := postgres.NewTransactionRepository(sqlxDB)
	pnlReportRepo := postgres.NewPnlReportRepository(sqlxDB)
	ChartOfAccountsRepo := postgres.NewChartOfAccountsRepository(sqlxDB)
	arrangementRepo := postgres.NewArrangementRepository(sqlxDB)

	// Use cases (application layer)
	clinicUC := usecase.NewClinicService(clinicRepo)
	userClinicUC := usecase.NewUserClinicService(userClinicRepo, clinicRepo, userRepo)
	authUC := usecase.NewAuthService(userRepo, sessionRepo, tokenService)
	expensesUC := usecase.NewExpensesService(expenseRepo)
	customFormUC := usecase.NewCustomFormService(customFormRepo, customFormFieldRepo, customFormVersionRepo, clinicRepo)
	clinicFinancialSettingsUC := usecase.NewClinicFinancialSettingsService(clinicFinancialSettingsRepo, clinicRepo, financialYearRepo, financialQuarterRepo, clinicFinancialYearRepo, clinicFinancialQuarterLockRepo)
	fieldEntryUC := usecase.NewFieldEntryService(fieldEntryRepo, customFormRepo, customFormFieldRepo, clinicRepo, clinicFinancialSettingsRepo, transactionRepo)
	transactionUC := usecase.NewTransactionService(transactionRepo, clinicRepo)
	pnlReportUC := usecase.NewPnlReportService(pnlReportRepo, clinicRepo)
	coaUC := usecase.NewCOAService(ChartOfAccountsRepo)
	arrangementUC := usecase.NewArrangementService(arrangementRepo)
	// HTTP handlers (driving adapters)
	authHandler := httpHandler.NewAuthHandler(authUC, oauthService, coaUC, cfg.OAuth.FrontendURL)
	userHandler := httpHandler.NewUserHandler(authUC, userClinicUC)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicUC, userClinicUC, clinicFinancialSettingsUC)
	userClinicHandler := httpHandler.NewUserClinicHandler(userClinicUC)
	customFormHandler := httpHandler.NewCustomFormHandler(customFormUC, userClinicUC)
	expensesHandler := httpHandler.NewExpensesHandler(expensesUC)
	clinicFinancialSettingsHandler := httpHandler.NewClinicFinancialSettingsHandler(clinicFinancialSettingsUC)
	fieldEntryHandler := httpHandler.NewFieldEntryHandler(fieldEntryUC, userClinicUC, customFormRepo)
	transactionHandler := httpHandler.NewTransactionHandler(transactionUC)
	pnlReportHandler := httpHandler.NewPnlReportHandler(pnlReportUC)
	coaHandler := httpHandler.NewCOAHandler(coaUC)
	arrangementHandler := httpHandler.NewArrangmentHandler(arrangementUC)

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
	custom_form.RegisterCustomFormRoutes(v1, customFormHandler, fieldEntryHandler, tokenService)
	expense.RegisterExpensesRoutes(v1, expensesHandler, tokenService)
	upload_route.RegisterUploadRoutes(v1, uploadHandler, tokenService)
	clinic_financial_settings.RegisterClinicFinancialSettingRoutes(v1, clinicFinancialSettingsHandler, tokenService)
	entry.RegisterEntryRoutes(v1, fieldEntryHandler, tokenService)
	transaction_route.RegisterTransactionRoutes(v1, transactionHandler, tokenService)
	reports_route.RegisterReportsRoutes(v1, pnlReportHandler, tokenService)
	coa_route.RegisterCOARoutes(v1, coaHandler, tokenService)
	arrangement.RegisterArrangementRoutes(v1, arrangementHandler, tokenService)
}
