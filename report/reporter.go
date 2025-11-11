package report

import (
	"fmt"
	"time"

	"go-load-lab/domain"
)

func toHuman(d time.Duration) string {
	if d > time.Second {
		return fmt.Sprintf("%.3fs", d.Seconds())
	}
	if d > time.Millisecond {
		return fmt.Sprintf("%.2fms", d.Seconds()*1000)
	}
	return fmt.Sprintf("%.2fµs", float64(d)/1000)
}

func Print(m domain.Metrics, duration time.Duration) {
	summary := domain.BuildSummary(m, duration)
	if summary.Latency == nil {
		fmt.Println("No successful requests")
		return
	}

	h := summary.Latency

	fmt.Printf("\n=== Load test results ===\n")
	fmt.Printf("Duration:      %v\n", duration)
	fmt.Printf("Total requests: %d (%.2f RPS)\n", summary.Total, summary.RPS)
	fmt.Printf("Success:        %d\n", summary.Success)
	fmt.Printf("Errors:         %d\n", summary.Errors)
	fmt.Printf("Latency:\n")
	fmt.Printf("  Min:    %s\n", toHuman(h.Min))
	fmt.Printf("  Max:    %s\n", toHuman(h.Max))
	fmt.Printf("  Mean:   %s\n", toHuman(h.Mean))
	fmt.Printf("  StdDev: %s\n", toHuman(h.StdDev))
	fmt.Printf("  P50:    %s\n", toHuman(h.P50))
	fmt.Printf("  P95:    %s\n", toHuman(h.P95))
	fmt.Printf("  P99:    %s\n", toHuman(h.P99))
	fmt.Printf("  P99.9:  %s\n", toHuman(h.P999))
	fmt.Printf("  P99.99: %s\n", toHuman(h.P9999))
}
