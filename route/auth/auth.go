package auth

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
)

func RegisterAuthRoutes(e *gin.RouterGroup, authHandler *httpHandler.AuthHandler) {
	auth := e.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/oauth/:provider", authHandler.InitiateOAuth)
	auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
}
