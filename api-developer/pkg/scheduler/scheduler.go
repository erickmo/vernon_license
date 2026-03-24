//go:build !wasm

// Package scheduler menyediakan scheduler untuk background jobs.
package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Job adalah fungsi yang dijalankan oleh scheduler.
type Job func(ctx context.Context) error

// Scheduler menjalankan job secara periodik.
type Scheduler struct {
	jobs   map[string]*scheduledJob
	log    *zap.Logger
	cancel context.CancelFunc
	done   chan struct{}
}

// scheduledJob menyimpan metadata satu scheduled job.
type scheduledJob struct {
	name     string
	interval time.Duration
	fn       Job
	ticker   *time.Ticker
	done     chan struct{}
}

// New membuat instance Scheduler baru.
func New(log *zap.Logger) *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*scheduledJob),
		log:  log,
		done: make(chan struct{}),
	}
}

// Schedule mendaftarkan job untuk dijalankan setiap interval.
// Job pertama kali dijalankan setelah interval, bukan langsung.
// Jika ada job dengan nama yang sama, akan di-replace.
func (s *Scheduler) Schedule(name string, interval time.Duration, fn Job) {
	s.jobs[name] = &scheduledJob{
		name:     name,
		interval: interval,
		fn:       fn,
		done:     make(chan struct{}),
	}
	s.log.Debug("Job scheduled", zap.String("name", name), zap.Duration("interval", interval))
}

// Start memulai semua scheduled jobs.
// Blocking call — hanya keluar saat ctx.Done().
func (s *Scheduler) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	for name, job := range s.jobs {
		go s.runJob(ctx, name, job)
	}

	<-ctx.Done()
	s.log.Info("Scheduler stopping")
	s.stopAllJobs()
	close(s.done)
}

// Stop menghentikan scheduler dan semua jobs-nya.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	<-s.done
	s.log.Info("Scheduler stopped")
}

// runJob menjalankan satu job secara periodik.
func (s *Scheduler) runJob(ctx context.Context, name string, job *scheduledJob) {
	job.ticker = time.NewTicker(job.interval)
	defer job.ticker.Stop()

	s.log.Info("Job started", zap.String("name", name))

	for {
		select {
		case <-job.ticker.C:
			s.executeJob(ctx, name, job)
		case <-job.done:
			return
		case <-ctx.Done():
			return
		}
	}
}

// executeJob menjalankan job function dengan error handling.
func (s *Scheduler) executeJob(ctx context.Context, name string, job *scheduledJob) {
	start := time.Now()
	err := job.fn(ctx)
	duration := time.Since(start)

	if err != nil {
		s.log.Error("Job failed",
			zap.String("name", name),
			zap.Duration("duration", duration),
			zap.Error(err))
		return
	}

	s.log.Debug("Job executed",
		zap.String("name", name),
		zap.Duration("duration", duration))
}

// stopAllJobs menghentikan semua running jobs.
func (s *Scheduler) stopAllJobs() {
	for name, job := range s.jobs {
		if job.ticker != nil {
			job.ticker.Stop()
		}
		close(job.done)
		s.log.Info("Job stopped", zap.String("name", name))
	}
}
