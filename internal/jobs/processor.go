package jobs

import (
	"context"

	"github.com/Sailesh-Dash/heliosflow/internal/logger"
	"github.com/Sailesh-Dash/heliosflow/internal/service"
)

// Processor pulls job IDs from a queue and processes them in
// the background using the JobService.
type Processor struct {
	svc   *service.JobService
	queue chan string
}

func NewProcessor(svc *service.JobService) *Processor {
	return &Processor{
		svc:   svc,
		queue: make(chan string, 100),
	}
}

// Enqueue schedules a job ID for background processing.
func (p *Processor) Enqueue(id string) {
	select {
	case p.queue <- id:
	default:
		logger.Error("job processor queue full, dropping job %s", id)
	}
}

// Start begins the background processing loop.
// Pass a long-lived context (e.g. context.Background()).
func (p *Processor) Start(ctx context.Context) {
	go func() {
		logger.Info("job processor started")
		defer logger.Info("job processor stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case id := <-p.queue:
				if err := p.svc.ProcessJob(ctx, id); err != nil {
					logger.Error("processing job %s failed: %v", id, err)
				}
			}
		}
	}()
}
