package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sailesh-Dash/heliosflow/internal/config"
	"github.com/Sailesh-Dash/heliosflow/internal/logger"
	httpserver "github.com/Sailesh-Dash/heliosflow/internal/server"
)

func main() {
	// Load configuration (e.g. HTTP port, etc.)
	cfg := config.FromEnv()

	// Create the HTTP server with all dependencies wired
	srv := httpserver.NewHTTPServer(cfg)

	// Listen for OS signals (Ctrl+C, kill, etc.)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in a background goroutine
	go func() {
		if err := srv.Start(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error("http server error: %v", err)
			os.Exit(1)
		}
	}()

	// Block until we get a shutdown signal
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown error: %v", err)
	}

	logger.Info("heliosflow API exited cleanly")
}
