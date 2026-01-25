package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bifshteksex/hertzboard/internal/config"
	"github.com/bifshteksex/hertzboard/internal/database"
	"github.com/bifshteksex/hertzboard/internal/handler"
	"github.com/bifshteksex/hertzboard/internal/repository"
	"github.com/bifshteksex/hertzboard/internal/router"
	"github.com/bifshteksex/hertzboard/internal/service"
	"github.com/cloudwego/hertz/pkg/app/server"
)

const (
	maxRequestBodySizeMB   = 10
	shutdownTimeoutSeconds = 5
	bytesInMB              = 1024 * 1024
	defaultConfigPath      = "configs/config.yaml"
)

func main() {
	// Initialize logger
	log.Println("Starting HertzBoard API Gateway...")

	// Load configuration
	configPath := getEnv("CONFIG_PATH", defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded configuration: %s environment", cfg.App.Env)

	// Connect to databases
	log.Println("Connecting to PostgreSQL...")
	dbPool, err := database.NewPostgresPool(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer database.ClosePostgresPool(dbPool)
	log.Println("Connected to PostgreSQL")

	log.Println("Connecting to Redis...")
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		database.ClosePostgresPool(dbPool)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		_ = database.CloseRedisClient(redisClient)
	}()
	log.Println("Connected to Redis")

	log.Println("Connecting to NATS...")
	natsConn, err := database.NewNATSConnection(&cfg.NATS)
	if err != nil {
		database.ClosePostgresPool(dbPool)
		_ = database.CloseRedisClient(redisClient)
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer database.CloseNATSConnection(natsConn)
	log.Println("Connected to NATS")

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.Migrate(dbPool, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)

	// Initialize services
	jwtService, err := service.NewJWTService(&cfg.JWT)
	if err != nil {
		log.Fatalf("Failed to create JWT service: %v", err)
	}

	authService := service.NewAuthService(userRepo, jwtService)
	oauthService := service.NewOAuthService(&cfg.OAuth, userRepo, jwtService)

	// Start email worker
	log.Println("Starting email worker...")
	emailWorker, err := service.NewEmailWorker(&cfg.Email, natsConn)
	if err != nil {
		log.Fatalf("Failed to start email worker: %v", err)
	}
	defer emailWorker.Close()
	log.Println("Email worker started")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo, authService)
	oauthHandler := handler.NewOAuthHandler(oauthService)

	// Initialize Hertz server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	h := server.Default(
		server.WithHostPorts(addr),
		server.WithMaxRequestBodySize(maxRequestBodySizeMB*bytesInMB),
	)

	// Setup routes and middleware
	deps := &router.Dependencies{
		JWTService:   jwtService,
		AuthHandler:  authHandler,
		UserHandler:  userHandler,
		OAuthHandler: oauthHandler,
	}
	router.Setup(h, cfg, deps)

	log.Printf("API Gateway is starting on %s", addr)

	// Graceful shutdown
	go func() {
		if err := h.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	log.Printf("API Gateway is running on %s", addr)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
