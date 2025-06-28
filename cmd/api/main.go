package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alxand/nalo-workspace/internal/container"
	"github.com/alxand/nalo-workspace/internal/pkg/logger"
	"github.com/alxand/nalo-workspace/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize container with all dependencies
	container, err := container.NewContainer()
	if err != nil {
		logger.Fatal("Failed to initialize container", zap.Error(err))
	}
	defer container.Close()

	// Create and configure the application
	app := server.NewApp(container.Config, container.Logger)
	app.SetupRoutes(
		container.AuthHandler,
		container.DailyTaskHandler,
		container.UserHandler,
		container.ContinentHandler,
		container.CountryHandler,
		container.CompanyHandler,
		container.AuthService,
	)

	// Start server in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			container.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		container.Logger.Error("Error during server shutdown", zap.Error(err))
	}

	container.Logger.Info("Server stopped")
}
