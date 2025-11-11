package concurrency

import (
	"context"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"go-load-lab/config"
	"go-load-lab/domain"
	"go-load-lab/infrastructure"
)

// WgRunner — реализация с sync.WaitGroup.
type WgRunner struct {
	Cfg config.Config
}

func (r *WgRunner) Run(
	ctx context.Context,
	executor *infrastructure.Executor,
	target domain.Target,
	results chan<- domain.Result,
	bar *progressbar.ProgressBar,
) {
	var wg sync.WaitGroup
	ticker := r.makeTicker()

	for i := 0; i < r.Cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if ticker != nil {
						<-ticker.C
					}
					res := executor.Do(target)
					results <- res
					_ = bar.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	if ticker != nil {
		ticker.Stop()
	}
}

func (r *WgRunner) makeTicker() *time.Ticker {
	if r.Cfg.Rate <= 0 {
		return nil
	}
	interval := time.Second / time.Duration(r.Cfg.Rate/r.Cfg.Concurrency+1)
	return time.NewTicker(interval)
}
