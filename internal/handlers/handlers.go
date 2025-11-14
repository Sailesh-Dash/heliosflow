package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Sailesh-Dash/heliosflow/internal/jobs"
	"github.com/Sailesh-Dash/heliosflow/internal/logger"
	"github.com/Sailesh-Dash/heliosflow/internal/service"
)

type Handlers struct {
	JobService *service.JobService
	Processor  *jobs.Processor
}

func NewHandlers(svc *service.JobService, processor *jobs.Processor) *Handlers {
	return &Handlers{
		JobService: svc,
		Processor:  processor,
	}
}

// Health is a simple liveness probe.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Ready is a simple readiness probe.
func (h *Handlers) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}

// Ping is a tiny test endpoint.
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

// CreateJob handles POST /v1/jobs.
func (h *Handlers) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		Name    string `json:"name"`
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	job, err := h.JobService.CreateJob(ctx, body.Name, body.Payload)
	if err != nil {
		logger.Error("create job failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Enqueue job for background processing.
	if h.Processor != nil {
		h.Processor.Enqueue(job.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(job)
}

// ListJobs handles GET /v1/jobs.
func (h *Handlers) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobsList, err := h.JobService.ListJobs(ctx)
	if err != nil {
		logger.Error("list jobs failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobsList)
}

// GetJob handles GET /v1/jobs/{id}.
func (h *Handlers) GetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	job, err := h.JobService.GetJob(ctx, id)
	if err != nil {
		if err == service.ErrJobNotFound {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		logger.Error("get job failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}

// CancelJob handles DELETE /v1/jobs/{id}.
func (h *Handlers) CancelJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.JobService.CancelJob(ctx, id)
	if err != nil {
		if err == service.ErrJobNotFound {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		logger.Error("cancel job failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
