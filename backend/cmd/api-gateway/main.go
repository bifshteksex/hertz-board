package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

const (
	defaultPort            = ":8080"
	maxRequestBodySizeMB   = 10
	shutdownTimeoutSeconds = 5
	bytesInMB              = 1024 * 1024
)

func main() {
	// Initialize logger
	log.Println("Starting HertzBoard API Gateway...")

	// Load configuration
	// TODO: Implement config loading

	// Initialize Hertz server
	h := server.Default(
		server.WithHostPorts(defaultPort),
		server.WithMaxRequestBodySize(maxRequestBodySizeMB*bytesInMB),
	)

	// TODO: Register routes
	// TODO: Register middleware
	// TODO: Connect to database
	// TODO: Connect to Redis
	// TODO: Initialize services

	// Register health check endpoint
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"service":   "api-gateway",
			"timestamp": time.Now().Unix(),
		})
	})

	// Graceful shutdown
	go func() {
		if err := h.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	log.Printf("API Gateway is running on %s", defaultPort)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		cancel() // Explicitly call cancel before Fatal
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
