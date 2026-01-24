package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// Initialize logger
	log.Println("Starting HertzBoard API Gateway...")

	// Load configuration
	// TODO: Implement config loading

	// Initialize Hertz server
	h := server.Default(
		server.WithHostPorts(":8080"),
		server.WithMaxRequestBodySize(10*1024*1024), // 10MB
	)

	// TODO: Register routes
	// TODO: Register middleware
	// TODO: Connect to database
	// TODO: Connect to Redis
	// TODO: Initialize services

	// Register health check endpoint
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(200, map[string]interface{}{
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

	log.Println("API Gateway is running on :8080")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited")
}
