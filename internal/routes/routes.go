package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Sailesh-Dash/heliosflow/internal/handlers"
	"github.com/Sailesh-Dash/heliosflow/internal/jobs"
	appmw "github.com/Sailesh-Dash/heliosflow/internal/middleware"
	"github.com/Sailesh-Dash/heliosflow/internal/repository"
	"github.com/Sailesh-Dash/heliosflow/internal/service"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Global middleware.
	r.Use(
		appmw.RequestID,
		appmw.Logging,
		appmw.Recoverer,
	)

	// Dependencies: repository -> service -> handlers.
	repo := repository.NewJobRepository()
	svc := service.NewJobService(repo)

	// Background worker for jobs.
	processor := jobs.NewProcessor(svc)
	processor.Start(context.Background())

	h := handlers.NewHandlers(svc, processor)

	// Basic health / readiness / ping.
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)
	r.Get("/v1/ping", h.Ping)

	// Job APIs.
	r.Route("/v1/jobs", func(r chi.Router) {
		r.Post("/", h.CreateJob)       // create a job
		r.Get("/", h.ListJobs)         // list jobs
		r.Get("/{id}", h.GetJob)       // get one job
		r.Delete("/{id}", h.CancelJob) // cancel a job
	})

	return r
}
