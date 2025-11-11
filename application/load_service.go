package application

import (
	"context"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"

	"go-load-lab/concurrency"
	"go-load-lab/config"
	"go-load-lab/domain"
	"go-load-lab/infrastructure"
)

type LoadService struct {
	runner concurrency.Runner
}

func NewLoadService(mode string, cfg config.Config) *LoadService {
	switch mode {
	case "channel":
		return &LoadService{runner: &concurrency.ChannelRunner{Cfg: cfg}}
	case "semaphore":
		return &LoadService{runner: &concurrency.SemaphoreRunner{Cfg: cfg}}
	default:
		return &LoadService{runner: &concurrency.WgRunner{Cfg: cfg}}
	}
}

func (s *LoadService) RunTest(cfg config.Config, target domain.Target, progress domain.Progress) domain.Metrics {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration+2*time.Second) // +2s на завершение
	defer cancel()

	if progress == nil {
		progress = domain.NewNoopProgress()
	}

	executor := infrastructure.NewExecutor(cfg)

	results := make(chan domain.Result, cfg.Concurrency*50)

	hist := hdrhistogram.New(1_000, 60_000_000_000, 3) // 1µs … 60s

	runnerDone := make(chan struct{})
	go func() {
		defer close(runnerDone)
		s.runner.Run(ctx, executor, target, results, progress)
	}()

	metrics := domain.Metrics{}
	collectorDone := make(chan struct{})

	go func() {
		defer close(collectorDone)
		for r := range results {
			metrics.Total++
			if r.Err == nil && r.Status >= 200 && r.Status < 300 {
				metrics.Success++
				_ = hist.RecordValue(r.Latency.Nanoseconds())
			} else {
				metrics.Errors++
			}
		}
	}()

	timer := time.NewTimer(cfg.Duration)
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	cancel()
	<-runnerDone
	close(results)
	<-collectorDone

	_ = progress.Finish()

	metrics.Hist = hist
	return metrics
}
