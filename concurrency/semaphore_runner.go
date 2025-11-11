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

type SemaphoreRunner struct {
	Cfg config.Config
}

func (r *SemaphoreRunner) Run(
	ctx context.Context,
	executor *infrastructure.Executor,
	target domain.Target,
	results chan<- domain.Result,
	bar *progressbar.ProgressBar,
) {
	sem := make(chan struct{}, r.Cfg.Concurrency)
	ticker := r.makeTicker()
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		default:
			if ticker != nil {
				<-ticker.C
			}

			sem <- struct{}{}
			wg.Add(1)

			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()

				res := executor.Do(target)
				results <- res
				_ = bar.Add(1)
			}()
		}
	}
}

func (r *SemaphoreRunner) makeTicker() *time.Ticker {
	if r.Cfg.Rate <= 0 {
		return nil
	}
	interval := time.Second / time.Duration(r.Cfg.Rate/r.Cfg.Concurrency+1)
	return time.NewTicker(interval)
}
