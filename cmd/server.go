package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/joho/godotenv"
)

func InitServer() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Run migrations
	if err := config.RunMigrations(db.DB.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	rdb := config.NewClient(cfg.Redis)
	defer rdb.Close()

	// Initialize services
	tokenService := service.NewTokenService(cfg.JWT)
	sessionService := service.NewSessionService(rdb, cfg.Session)
	oauthService := service.NewOAuthService(cfg.OAuth, db.DB)
	authService := service.NewAuthService(db.DB, tokenService, sessionService, oauthService)

	// Initialize handlers
	authHandler := domain.NewAuthHandler(authService, oauthService)
	userHandler := domain.NewUserHandler(authService)

	// Setup router
	router := setupRouter(authHandler, userHandler, tokenService)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("SSO service started on port %s", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(authHandler *domain.AuthHandler, userHandler *domain.UserHandler, tokenService *service.TokenService) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public routes
	api := router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(tokenService), authHandler.Logout)

			// OAuth routes
			auth.GET("/oauth/:provider", authHandler.InitiateOAuth)
			auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
		}

		// Protected routes
		protected := api.Group("/users")
		protected.Use(middleware.AuthMiddleware(tokenService))
		{
			protected.GET("/me", userHandler.GetCurrentUser)
			protected.PUT("/me", userHandler.UpdateCurrentUser)
			protected.GET("/sessions", userHandler.GetActiveSessions)
			protected.DELETE("/sessions/:sessionId", userHandler.RevokeSession)
		}
	}

	return router
}
