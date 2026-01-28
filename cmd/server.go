package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/iamarpitzala/aca-reca-backend/route"
	"github.com/joho/godotenv"
)

func InitServer() {
	// Load environment variables
	var envOnce sync.Once
	envOnce.Do(func() {
		envFile := ".env"

		// Try loading the .env file
		err := godotenv.Load(envFile)
		if err != nil {
			// If .env file is not found, don't log an error
			if !os.IsNotExist(err) {
				log.Fatalf("Error loading .env file from %s: %s", envFile, err)
			} else {
				log.Println(".env file not found, falling back to system environment variables")
			}
		} else {
			log.Println(".env file loaded successfully")
		}
	})

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

	e := gin.New()

	// Configure CORS
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://preview--zenithive.lovable.app", "https://zenithive.lovable.app", "https://zenithive.lovable.app/", "https://preview--zenithive.lovable.app/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	e.Use(gin.Recovery())
	e.Use(gin.Logger())

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	route.InitRouter(e)

	port := cfg.Server.Port

	if port == "" {
		port = "8080"
	}

	//Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	// // Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("ACA RECA service started on port %s", port)

	// // Wait for interrupt signal
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
