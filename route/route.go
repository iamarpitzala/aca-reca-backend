package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	_ "github.com/iamarpitzala/aca-reca-backend/docs" // swagger docs
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/iamarpitzala/aca-reca-backend/route/aoc"
	"github.com/iamarpitzala/aca-reca-backend/route/auth"
	"github.com/iamarpitzala/aca-reca-backend/route/clinic"
	expense "github.com/iamarpitzala/aca-reca-backend/route/expense"
	financial_form "github.com/iamarpitzala/aca-reca-backend/route/financial_form"
	payslip "github.com/iamarpitzala/aca-reca-backend/route/payship"
	"github.com/iamarpitzala/aca-reca-backend/route/quarter"
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

	tokenService := service.NewTokenService(cfg.JWT)
	oauthService := service.NewOAuthService(cfg.OAuth, db.DB)
	authService := service.NewAuthService(db.DB, tokenService, oauthService)
	clinicService := service.NewClinicService(db.DB)
	userClinicService := service.NewUserClinicService(db.DB)
	financialFormService := service.NewFinancialFormService(db.DB)
	// financialCalculationService := service.NewFinancialCalculationService(db.DB)
	expensesService := service.NewExpensesService(db.DB)
	quarterService := service.NewQuarterService(db.DB)
	aosService := service.NewAOSService(db.DB)

	authHandler := httpHandler.NewAuthHandler(authService, oauthService, cfg.OAuth.FrontendURL)
	userHandler := httpHandler.NewUserHandler(authService)
	payslipHandler := httpHandler.NewPayslipHandler()
	clinicHandler := httpHandler.NewClinicHandler(clinicService, userClinicService)
	userClinicHandler := httpHandler.NewUserClinicHandler(userClinicService)
	financialFormHandler := httpHandler.NewFinancialFormHandler(financialFormService)
	// financialCalculationHandler := httpHandler.NewFinancialCalculationHandler(financialCalculationService)
	expensesHandler := httpHandler.NewExpensesHandler(expensesService)
	quarterHandler := httpHandler.NewQuarterHandler(quarterService)
	aosHandler := httpHandler.NewAOCHandler(aosService)
	// Swagger documentation route
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/api/v1")
	auth.RegisterAuthRoutes(v1, authHandler)
	auth.RegisterUserRoutes(v1, userHandler)
	clinic.RegisterClinicRoutes(v1, clinicHandler)
	payslip.RegisterPayslipRoutes(v1, payslipHandler)
	user_clinic.RegisterUserClinicRoutes(v1, userClinicHandler)
	financial_form.RegisterFinancialFormRoutes(v1, financialFormHandler)
	// financial_calculation.RegisterFinancialCalculationRoutes(v1, financialCalculationHandler)
	quarter.RegisterQuarterRoutes(v1, quarterHandler)
	expense.RegisterExpensesRoutes(v1, expensesHandler)
	aoc.RegisterAOCRoutes(v1, aosHandler)

}
