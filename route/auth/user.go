package auth

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
)

func RegisterUserRoutes(e *gin.RouterGroup, userHandler *httpHandler.UserHandler) {
	user := e.Group("/users")
	{
		user.GET("/:userId", userHandler.GetCurrentUser)
		user.PUT("/:userId", userHandler.UpdateCurrentUser)
		user.GET("/:userId/sessions", userHandler.GetActiveSessions)
		user.DELETE("/sessions/:sessionId", userHandler.RevokeSession)
		user.POST("/:userId/sessions/revoke", userHandler.RevokeSession)
	}
}
