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
	expense.GET("/type/clinic/:clinicId", expensesHandler.GetExpenseTypesByClinicID)
	expense.GET("/type/:id", expensesHandler.GetExpenseTypeByID)
	expense.GET("/category/clinic/:clinicId", expensesHandler.GetExpenseCategoriesByClinicID)
	expense.GET("/category/:id", expensesHandler.GetExpenseCategoryByID)
	expense.PUT("/category/:id", expensesHandler.UpdateExpenseCategory)
	expense.DELETE("/category/:id", expensesHandler.DeleteExpenseCategory)
	expense.GET("/category-type/:id", expensesHandler.GetExpenseCategoryTypeByID)
	expense.GET("/entry/clinic/:clinicId", expensesHandler.GetExpenseEntriesByClinicID)
	expense.GET("/entry/:id", expensesHandler.GetExpenseEntryByID)
}
