package auth

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterUserRoutes(e *gin.RouterGroup, userHandler *httpHandler.UserHandler, tokenService *service.TokenService) {
	user := e.Group("/users")
	user.Use(middleware.AuthMiddleware(tokenService))

	// /me must be before /:userId to avoid "me" being captured as userId
	user.GET("/me", userHandler.GetMe)
	user.GET("/:userId", userHandler.GetCurrentUser)
	user.PUT("/:userId", userHandler.UpdateCurrentUser)
	user.GET("/:userId/sessions", userHandler.GetActiveSessions)
	user.DELETE("/:userId/sessions/:sessionId", userHandler.RevokeSession)
}
