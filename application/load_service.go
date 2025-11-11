package application

import (
	"context"
	"math"
	"os"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/schollz/progressbar/v3"

	"go-load-lab/config"
	"go-load-lab/concurrency"
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


func (s *LoadService) RunTest(cfg config.Config, target domain.Target) domain.Metrics {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration+2*time.Second) // +2s на завершение
	defer cancel()

	executor := infrastructure.NewExecutor(cfg)

	results := make(chan domain.Result, cfg.Concurrency*50)

	// progress-bar
	totalExpected := int64(-1)
	if cfg.Rate > 0 {
		totalExpected = int64(math.Ceil(float64(cfg.Rate) * cfg.Duration.Seconds())) + int64(cfg.Concurrency)*2
	}
	bar := progressbar.NewOptions64(totalExpected,
		progressbar.OptionSetDescription("Load test"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(totalExpected > 0),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionSpinnerType(14),
	)

	hist := hdrhistogram.New(1_000, 60_000_000_000, 3) // 1µs … 60s

	go s.runner.Run(ctx, executor, target, results, bar)

	metrics := domain.Metrics{}
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(done)
				return
			case r := <-results:
				metrics.Total++
				if r.Err == nil && r.Status >= 200 && r.Status < 300 {
					metrics.Success++
					_ = hist.RecordValue(r.Latency.Nanoseconds())
				} else {
					metrics.Errors++
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
	case <-time.After(cfg.Duration):
	}

	bar.Finish() // ждём пока все запросы завершатся
	<-done
	close(results)
	bar.Finish()

	metrics.Hist = hist
	return metrics
}