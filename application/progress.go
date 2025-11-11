package application

import (
	"math"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"

	"go-load-lab/config"
	"go-load-lab/domain"
)

// NewProgress creates a progress reporter suitable for interactive CLI usage.
// When rate is provided we set an expected total to improve ETA predictions.
func NewProgress(cfg config.Config) domain.Progress {
	totalExpected := int64(-1)
	if cfg.Rate > 0 {
		totalExpected = int64(math.Ceil(float64(cfg.Rate)*cfg.Duration.Seconds())) + int64(cfg.Concurrency)*2
	}

	return progressbar.NewOptions64(totalExpected,
		progressbar.OptionSetDescription("Load test"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(totalExpected > 0),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionSpinnerType(14),
	)
}
