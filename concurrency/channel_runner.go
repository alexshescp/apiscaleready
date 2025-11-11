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

// ChannelRunner — реализация через каналы и воркеры.
type ChannelRunner struct {
	Cfg config.Config // Сделали экспортируемым (с заглавной)
}

func (r *ChannelRunner) Run(
	ctx context.Context,
	executor *infrastructure.Executor,
	target domain.Target,
	results chan<- domain.Result,
	bar *progressbar.ProgressBar,
) {
	jobs := make(chan struct{}, r.Cfg.Concurrency*10)

	// Rate limiter
	if r.Cfg.Rate > 0 {
		go func() {
			ticker := time.NewTicker(time.Second / time.Duration(r.Cfg.Rate))
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					jobs <- struct{}{}
				}
			}
		}()
	} else {
		for i := 0; i < cap(jobs); i++ {
			jobs <- struct{}{}
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < r.Cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-jobs:
					res := executor.Do(target)
					results <- res
					_ = bar.Add(1)
				}
			}
		}()
	}
	wg.Wait()
}
