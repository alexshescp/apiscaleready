package domain

import (
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

type Target struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
}

type Result struct {
	Latency time.Duration
	Status  int
	Err     error
}

type Metrics struct {
	Total   uint64
	Success uint64
	Errors  uint64
	Hist    *hdrhistogram.Histogram // правильный тип из пакета hdrhistogram-go
}
