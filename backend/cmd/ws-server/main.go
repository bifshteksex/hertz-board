package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// Initialize logger
	log.Println("Starting HertzBoard WebSocket Server...")

	// Load configuration
	// TODO: Implement config loading

	// Initialize Hertz server for WebSocket
	h := server.Default(
		server.WithHostPorts(":8081"),
	)

	// TODO: Initialize WebSocket hub
	// TODO: Register WebSocket handlers
	// TODO: Connect to Redis for pub/sub
	// TODO: Initialize CRDT sync engine

	// Register health check endpoint
	h.GET("/health", func(c context.Context, ctx *server.RequestContext) {
		ctx.JSON(200, map[string]interface{}{
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

	log.Println("WebSocket Server is running on :8081")

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
