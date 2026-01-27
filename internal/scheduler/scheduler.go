package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/nbaisland/nbaisland/internal/logger"
)

type Job struct {
	Name     string
	Schedule time.Duration
	RunAt    time.Time
	Fn       func(ctx context.Context) error
}

type Scheduler struct {
	jobs []Job
}

func New() *Scheduler {
	return &Scheduler{
		jobs: make([]Job, 0),
	}
}

func (s *Scheduler) AddWeekly(name string, runAtHour, runAtMinute int, fn func(ctx context.Context) error) {
	s.jobs = append(s.jobs, Job{
		Name:     name,
		Schedule: 7 * 24 * time.Hour,
		RunAt:    time.Date(0, 1, 1, runAtHour, runAtMinute, 0, 0, time.UTC),
		Fn:       fn,
	})
}

func (s *Scheduler) AddNightly(name string, runAtHour, runAtMinute int, fn func(ctx context.Context) error) {
	s.jobs = append(s.jobs, Job{
		Name:     name,
		Schedule: 24 * time.Hour,
		RunAt:    time.Date(0, 1, 1, runAtHour, runAtMinute, 0, 0, time.UTC),
		Fn:       fn,
	})
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		go s.runJob(ctx, job)
	}
}

func (s *Scheduler) runJob(ctx context.Context, job Job) {
	nextRun := s.calculateNextRun(job.RunAt, job.Schedule)

	logger.Log.Info(
		"scheduled job initialized",
		zap.String("job", job.Name),
		zap.Time("first_run_at", nextRun),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info(
				"scheduled job stopped",
				zap.String("job", job.Name),
			)
			return

		case <-time.After(time.Until(nextRun)):
			logger.Log.Info(
				"scheduled job starting",
				zap.String("job", job.Name),
				zap.Time("run_at", nextRun),
			)

			if err := job.Fn(ctx); err != nil {
				logger.Log.Error(
					"scheduled job failed",
					zap.String("job", job.Name),
					zap.Error(err),
				)
			}

			nextRun = nextRun.Add(job.Schedule)

			logger.Log.Info(
				"scheduled job completed",
				zap.String("job", job.Name),
				zap.Time("next_run_at", nextRun),
			)
		}
	}
}

func (s *Scheduler) calculateNextRun(runAt time.Time, schedule time.Duration) time.Time {
	now := time.Now()

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		runAt.Hour(),
		runAt.Minute(),
		0,
		0,
		now.Location(),
	)

	if now.After(next) {
		next = next.Add(schedule)
	}

	return next
}
