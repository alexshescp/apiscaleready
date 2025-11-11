package main

import (
	"log"
	"time"

	"go-load-lab/application"
	"go-load-lab/config"
	"go-load-lab/domain"
	"go-load-lab/presentation/httpserver"
	"go-load-lab/report"
)

func main() {
	cfg := config.Load()

	if cfg.ListenAddr != "" {
		srv := httpserver.New(cfg)
		log.Fatal(srv.Start())
		return
	}

	if cfg.TargetURL == "" {
		log.Fatal("GLAB_TARGET_URL is required when HTTP API is disabled")
	}

	target := domain.Target{
		URL:     cfg.TargetURL,
		Method:  cfg.Method,
		Headers: cfg.Headers,
		Body:    []byte(cfg.JsonBody),
	}

	service := application.NewLoadService(cfg.Mode, cfg)
	progress := application.NewProgress(cfg)

	start := time.Now()
	metrics := service.RunTest(cfg, target, progress)
	duration := time.Since(start)

	report.Print(metrics, duration)
}
