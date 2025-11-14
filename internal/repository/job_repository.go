package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/Sailesh-Dash/heliosflow/internal/models"
)

var ErrJobNotFound = errors.New("job not found")

// JobRepository is an in-memory store for jobs.
type JobRepository struct {
	mu   sync.RWMutex
	jobs map[string]models.Job
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs: make(map[string]models.Job),
	}
}

func (r *JobRepository) Create(job models.Job) (models.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID] = job
	return job, nil
}

func (r *JobRepository) GetByID(id string) (models.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		return models.Job{}, ErrJobNotFound
	}
	return job, nil
}

func (r *JobRepository) List() ([]models.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]models.Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		out = append(out, j)
	}
	return out, nil
}

// UpdateStatus updates the status (and UpdatedAt) of a job.
func (r *JobRepository) UpdateStatus(id string, status models.JobStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[id]
	if !ok {
		return ErrJobNotFound
	}

	job.Status = status
	job.UpdatedAt = time.Now().UTC()
	r.jobs[id] = job

	return nil
}

// Cancel marks a job as canceled.
func (r *JobRepository) Cancel(id string) error {
	return r.UpdateStatus(id, models.StatusCanceled)
}
