package service

import (
	"context"
	"time"

	"github.com/Sailesh-Dash/heliosflow/internal/models"
	"github.com/Sailesh-Dash/heliosflow/internal/repository"
)

var ErrJobNotFound = repository.ErrJobNotFound

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (s *JobService) CreateJob(ctx context.Context, name, payload string) (models.Job, error) {
	job := models.NewJob(name, payload)
	return s.repo.Create(job)
}

func (s *JobService) GetJob(ctx context.Context, id string) (models.Job, error) {
	return s.repo.GetByID(id)
}

func (s *JobService) ListJobs(ctx context.Context) ([]models.Job, error) {
	return s.repo.List()
}

func (s *JobService) CancelJob(ctx context.Context, id string) error {
	return s.repo.Cancel(id)
}

// ProcessJob simulates background processing of a job.
func (s *JobService) ProcessJob(ctx context.Context, id string) error {
	// Mark as running
	if err := s.repo.UpdateStatus(id, models.StatusRunning); err != nil {
		return err
	}

	// Simulate some work
	select {
	case <-ctx.Done():
		// If we were shut down, mark as failed and return
		_ = s.repo.UpdateStatus(id, models.StatusFailed)
		return ctx.Err()
	case <-time.After(5 * time.Second):
	}

	// Mark as success
	if err := s.repo.UpdateStatus(id, models.StatusSuccess); err != nil {
		return err
	}

	return nil
}
