package financial_form

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterFinancialFormRoutes(e *gin.RouterGroup, financialFormHandler *httpHandler.FinancialFormHandler) {
	financialForm := e.Group("/financial-form")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	financialForm.Use(middleware.AuthMiddleware(tokenService))

	financialForm.POST("/", financialFormHandler.CreateFinancialForm)
	financialForm.GET("/:id", financialFormHandler.GetFinancialForm)
	financialForm.GET("/clinic/:clinicId", financialFormHandler.GetFinancialFormsByClinic)
	financialForm.PUT("/:id", financialFormHandler.UpdateFinancialForm)
	financialForm.DELETE("/:id", financialFormHandler.DeleteFinancialForm)
}
