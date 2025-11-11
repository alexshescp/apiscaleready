package concurrency

import (
	"context"

	"github.com/schollz/progressbar/v3"
	"go-load-lab/domain"
	"go-load-lab/infrastructure"
)

// Runner — интерфейс для разных стратегий запуска нагрузки.
type Runner interface {
	Run(
		ctx context.Context,
		executor *infrastructure.Executor,
		target domain.Target,
		results chan<- domain.Result,
		bar *progressbar.ProgressBar,
	)
}
