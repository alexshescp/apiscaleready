package domain

import "time"

// LatencySummary describes aggregated latency percentiles.
type LatencySummary struct {
        Min    time.Duration `json:"min"`
        Max    time.Duration `json:"max"`
        Mean   time.Duration `json:"mean"`
        StdDev time.Duration `json:"std_dev"`
        P50    time.Duration `json:"p50"`
        P95    time.Duration `json:"p95"`
        P99    time.Duration `json:"p99"`
        P999   time.Duration `json:"p999"`
        P9999 time.Duration `json:"p9999"`
}

// Summary is a structured representation of Metrics that is convenient for JSON
// serialisation and documentation.
type Summary struct {
        Total   uint64          `json:"total"`
        Success uint64          `json:"success"`
        Errors  uint64          `json:"errors"`
        RPS     float64         `json:"rps"`
        Latency *LatencySummary `json:"latency,omitempty"`
}

// BuildSummary converts Metrics into a Summary using the provided total test
// duration.
func BuildSummary(m Metrics, duration time.Duration) Summary {
        summary := Summary{
                Total:   m.Total,
                Success: m.Success,
                Errors:  m.Errors,
        }

        if duration > 0 {
                summary.RPS = float64(m.Total) / duration.Seconds()
        }

        if m.Hist != nil && m.Hist.TotalCount() > 0 {
                summary.Latency = &LatencySummary{
                        Min:    time.Duration(m.Hist.Min()),
                        Max:    time.Duration(m.Hist.Max()),
                        Mean:   time.Duration(int64(m.Hist.Mean())),
                        StdDev: time.Duration(int64(m.Hist.StdDev())),
                        P50:    time.Duration(m.Hist.ValueAtPercentile(50)),
                        P95:    time.Duration(m.Hist.ValueAtPercentile(95)),
                        P99:    time.Duration(m.Hist.ValueAtPercentile(99)),
                        P999:   time.Duration(m.Hist.ValueAtPercentile(99.9)),
                        P9999:  time.Duration(m.Hist.ValueAtPercentile(99.99)),
                }
        }

        return summary
}
