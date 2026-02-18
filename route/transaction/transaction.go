package transaction

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterTransactionRoutes(e *gin.RouterGroup, txnHandler *httpHandler.TransactionHandler, tokenService *service.TokenService) {
	txn := e.Group("/transaction")
	txn.Use(middleware.AuthMiddleware(tokenService))

	txn.POST("", txnHandler.Create)
	txn.GET("/:id", txnHandler.GetByID)
	txn.GET("/clinic/:clinicId", txnHandler.GetByClinicID)
	txn.PUT("/:id", txnHandler.Update)
	txn.POST("/:id/post", txnHandler.Post)
	txn.POST("/:id/void", txnHandler.Void)
	txn.DELETE("/:id", txnHandler.Delete)
}
