package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/marcell322/api-test-studio/internal/adapters/persistence"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/server"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

func main() {
	// load .env file (optional, won't fail if missing)
	if err := godotenv.Load(); err != nil {
		log.Println("note: .env file not found, using environment variables")
	}

	// load configuration from environment
	cfg := config.Load()

	// initialize database connection
	log.Println("initializing database...")
	db, err := persistence.NewGormDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	log.Printf("database initialized at %s", cfg.DBPath)

	// auto migrate models
	log.Println("running migrations...")
	if err := persistence.AutoMigrate(db, &models.User{}); err != nil {
		log.Fatalf("automigrate failed: %v", err)
	}
	log.Println("migrations completed")

	// initialize repositories and services
	log.Println("initializing services...")
	userRepo := persistence.NewGormUserRepository(db)
	userSvc := usecase.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpireH)
	log.Println("services initialized")

	// setup router and register routes
	log.Println("configuring routes...")
	r := server.NewRouter(cfg, userSvc)
	log.Println("routes configured")

	// create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start server in a goroutine
	go func() {
		log.Printf("server listening on %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown: listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// close database connection
	if err := persistence.CloseDB(db); err != nil {
		log.Printf("database close error: %v", err)
	}

	log.Println("server stopped")
}
