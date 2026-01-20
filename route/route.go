package route

import (
	payslip "github.com/iamarpitzala/aca-reca-backend/route/payship"

	_ "github.com/iamarpitzala/aca-reca-backend/docs"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine) {
	//cfg := config.Load()
	//db, err := config.NewConnection(cfg.DB)
	//if err != nil {
	//log.Fatalf("Failed to connect to database: %v", err)
	//}

	rdb := config.NewRedisClient(cfg.Redis)

	// tokenService := service.NewTokenService(cfg.JWT)
	// sessionService := service.NewSessionService(rdb, cfg.Session)
	// oauthService := service.NewOAuthService(cfg.OAuth, db.DB)
	// authService := service.NewAuthService(db.DB, tokenService, sessionService, oauthService)

	// authHandler := httpHandler.NewAuthHandler(authService, oauthService)
	// userHandler := httpHandler.NewUserHandler(authService)
	payslipHandler := httpHandler.NewPayshipHandler()

	// Swagger documentation route
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/api/v1")
	// auth.RegisterAuthRoutes(v1, authHandler)
	// auth.RegisterUserRoutes(v1, userHandler)
	payslip.RegisterPayshipRoutes(v1, payslipHandler)

}

// // applySubscriptionMiddleware applies subscription middleware to protected routes
// func applySubscriptionMiddleware(api *gin.RouterGroup, subMiddleware *middleware.SubscriptionMiddleware) {
// 	// Apply optional subscription middleware to all API routes
// 	// This adds subscription info to context if available
// 	api.Use(subMiddleware.OptionalSubscription())

// 	// Protected routes that require active subscription
// 	protected := api.Group("/protected")
// 	{
// 		// Routes that require active subscription
// 		protected.Use(subMiddleware.RequireActiveSubscription())
// 		// Add protected routes here
// 	}

// 	// Premium routes that require premium subscription
// 	premium := api.Group("/premium")
// 	{
// 		// Routes that require premium subscription
// 		premium.Use(subMiddleware.RequirePremiumSubscription())
// 		// Add premium routes here
// 	}
// }
