package scheduler

import (
    "context"
    "log"
    "time"
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

func (s *Scheduler) Start(ctx context.Context) {
    for _, job := range s.jobs {
        go s.runJob(ctx, job)
    }
}

func (s *Scheduler) runJob(ctx context.Context, job Job) {
    nextRun := s.calculateNextRun(job.RunAt)
    
    log.Printf("Scheduled job '%s' will first run at %s", job.Name, nextRun.Format(time.RFC3339))
    
    for {
        select {
        case <-ctx.Done():
            log.Printf("Stopping scheduled job: %s", job.Name)
            return
        case <-time.After(time.Until(nextRun)):
            log.Printf("Running scheduled job: %s", job.Name)
            if err := job.Fn(ctx); err != nil {
                log.Printf("Error in job %s: %v", job.Name, err)
            }
            
            nextRun = nextRun.Add(job.Schedule)
            log.Printf("Next run for '%s' scheduled at %s", job.Name, nextRun.Format(time.RFC3339))
        }
    }
}

func (s *Scheduler) calculateNextRun(runAt time.Time) time.Time {
    now := time.Now()
    
    next := time.Date(
        now.Year(),
        now.Month(),
        now.Day(),
        runAt.Hour(),
        runAt.Minute(),
        0, 0,
        now.Location(),
    )
    
    if now.After(next) {
        next = next.Add(7 * 24 * time.Hour)
    }
    
    return next
}