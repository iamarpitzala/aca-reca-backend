package auth

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
)

func RegisterUserRoutes(e *gin.RouterGroup, userHandler *httpHandler.UserHandler) {
	user := e.Group("/users")
	{
		user.GET("/:user_id", userHandler.GetCurrentUser)
		user.PUT("/:user_id", userHandler.UpdateCurrentUser)
		user.GET("/sessions", userHandler.GetActiveSessions)
		user.DELETE("/sessions/:session_id", userHandler.RevokeSession)
		user.POST("/sessions/revoke", userHandler.RevokeSession)
	}
}
