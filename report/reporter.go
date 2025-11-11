package report

import (
	"fmt"
	"time"
	"go-load-lab/domain"
)

func toHuman(ns int64) string {
	d := time.Duration(ns) * time.Nanosecond
	if d > time.Second {
		return fmt.Sprintf("%.3fs", d.Seconds())
	}
	if d > time.Millisecond {
		return fmt.Sprintf("%.2fms", d.Seconds()*1000)
	}
	return fmt.Sprintf("%.2fµs", float64(d)/1000)
}

func Print(m domain.Metrics, duration time.Duration) {
	if m.Hist == nil || m.Hist.TotalCount() == 0 {
		fmt.Println("No successful requests")
		return
	}

	h := m.Hist
	rps := float64(m.Total) / duration.Seconds()

	fmt.Printf("\n=== Load test results ===\n")
	fmt.Printf("Duration:      %v\n", duration)
	fmt.Printf("Total requests: %d (%.2f RPS)\n", m.Total, rps)
	fmt.Printf("Success:        %d\n", m.Success)
	fmt.Printf("Errors:         %d\n", m.Errors)
	fmt.Printf("Latency:\n")
	fmt.Printf("  Min:    %s\n", toHuman(h.Min()))
	fmt.Printf("  Max:    %s\n", toHuman(h.Max()))
	fmt.Printf("  Mean:   %s\n", toHuman(int64(h.Mean())))
	fmt.Printf("  StdDev: %s\n", toHuman(int64(h.StdDev())))
	fmt.Printf("  P50:    %s\n", toHuman(h.ValueAtPercentile(50)))
	fmt.Printf("  P95:    %s\n", toHuman(h.ValueAtPercentile(95)))
	fmt.Printf("  P99:    %s\n", toHuman(h.ValueAtPercentile(99)))
	fmt.Printf("  P99.9:  %s\n", toHuman(h.ValueAtPercentile(99.9)))
	fmt.Printf("  P99.99: %s\n", toHuman(h.ValueAtPercentile(99.99)))
}