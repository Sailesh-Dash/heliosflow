package models

import (
	"time"

	"github.com/google/uuid"
)

// JobStatus is a strongly-typed status for jobs.
type JobStatus string

const (
	StatusPending  JobStatus = "pending"
	StatusRunning  JobStatus = "running"
	StatusSuccess  JobStatus = "success"
	StatusFailed   JobStatus = "failed"
	StatusCanceled JobStatus = "canceled"
)

// Job represents a background job in the system.
type Job struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Payload   string    `json:"payload,omitempty"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

// NewJob constructs a new Job in pending state.
func NewJob(name, payload string) Job {
	now := time.Now().UTC()

	return Job{
		ID:        uuid.NewString(),
		Name:      name,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
