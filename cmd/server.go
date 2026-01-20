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
	"github.com/iamarpitzala/aca-reca-backend/route"
)

func InitServer() {
	// Load environment variables
	// if err := godotenv.Load("./.env"); err != nil {
	// 	log.Println("No .env file found, using environment variables")
	// }

	// // Load configuration
	// cfg := config.Load()

	// // Initialize database
	// db, err := config.NewConnection(cfg.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer func() {
	// 	if err := db.Close(); err != nil {
	// 		log.Printf("Error closing database: %v", err)
	// 	}
	// }()

	// // Run migrations
	// if err := config.RunMigrations(db.DB.DB); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// // Initialize Redis
	// rdb := config.NewRedisClient(cfg.Redis)
	// defer rdb.Close()

	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())

	route.InitRouter(e)

	//port := cfg.Server.Port
	port := "8081"
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

	// log.Printf("ACA RECA service started on port %s", port)

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
