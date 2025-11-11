package concurrency

import (
	"context"

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
		progress domain.Progress,
	)
}
