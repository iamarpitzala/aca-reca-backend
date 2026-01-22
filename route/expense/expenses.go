package expense

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterExpensesRoutes(e *gin.RouterGroup, expensesHandler *httpHandler.ExpensesHandler) {
	expense := e.Group("/expense")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	expense.Use(middleware.AuthMiddleware(tokenService))

	expense.POST("/type", expensesHandler.CreateExpenseType)
	expense.POST("/category", expensesHandler.CreateExpenseCategory)
	expense.POST("/category-type", expensesHandler.CreateExpenseCategoryType)
	expense.POST("/entry", expensesHandler.CreateExpenseEntry)
	expense.GET("/type/:id", expensesHandler.GetExpenseTypeByID)
	expense.GET("/category/:id", expensesHandler.GetExpenseCategoryByID)
	expense.GET("/category-type/:id", expensesHandler.GetExpenseCategoryTypeByID)
	expense.GET("/entry/:id", expensesHandler.GetExpenseEntryByID)
}
