package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Sailesh-Dash/heliosflow/internal/config"
	"github.com/Sailesh-Dash/heliosflow/internal/logger"
	"github.com/Sailesh-Dash/heliosflow/internal/routes"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(cfg config.Config) *HTTPServer {
	router := routes.RegisterRoutes()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &HTTPServer{srv: srv}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	logger.Info("HTTP server starting on %s", s.srv.Addr)

	// Start HTTP server
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	logger.Info("HTTP server shutting down")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.srv.Shutdown(shutdownCtx)
}
