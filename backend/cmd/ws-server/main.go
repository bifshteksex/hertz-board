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
	defaultPort            = ":8081"
	shutdownTimeoutSeconds = 5
)

func main() {
	// Initialize logger
	log.Println("Starting HertzBoard WebSocket Server...")

	// Load configuration
	// TODO: Implement config loading

	// Initialize Hertz server for WebSocket
	h := server.Default(
		server.WithHostPorts(defaultPort),
	)

	// TODO: Initialize WebSocket hub
	// TODO: Register WebSocket handlers
	// TODO: Connect to Redis for pub/sub
	// TODO: Initialize CRDT sync engine

	// Register health check endpoint
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"service":   "ws-server",
			"timestamp": time.Now().Unix(),
		})
	})

	// Graceful shutdown
	go func() {
		if err := h.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	log.Printf("WebSocket Server is running on %s", defaultPort)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)

	if err := h.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	cancel()
	fmt.Println("Server exited")
}
