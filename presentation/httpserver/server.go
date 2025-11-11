package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"go-load-lab/application"
	"go-load-lab/config"
	"go-load-lab/domain"
)

// Server exposes a minimal HTTP API for triggering load tests.
type Server struct {
	cfg config.Config
}

// New creates a new Server instance.
func New(cfg config.Config) *Server {
	return &Server{cfg: cfg}
}

// Handler returns the HTTP handler exposing the API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/tests", s.handleRunTest)
	return mux
}

// Start runs the HTTP server and blocks until it stops.
func (s *Server) Start() error {
	log.Printf("HTTP API listening on %s", s.cfg.ListenAddr)
	return http.ListenAndServe(s.cfg.ListenAddr, s.Handler())
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type testRequest struct {
	TargetURL   string            `json:"target_url"`
	Method      string            `json:"method"`
	Concurrency *int              `json:"concurrency"`
	Duration    string            `json:"duration"`
	Rate        *int              `json:"rate"`
	Mode        string            `json:"mode"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	NoKeepAlive *bool             `json:"no_keep_alive"`
}

type testResponse struct {
	Duration string         `json:"duration"`
	Metrics  domain.Summary `json:"metrics"`
}

func (s *Server) handleRunTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req testRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cfg, target, err := s.buildConfig(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := application.NewLoadService(cfg.Mode, cfg)

	start := time.Now()
	metrics := service.RunTest(cfg, target, domain.NewNoopProgress())
	duration := time.Since(start)

	resp := testResponse{
		Duration: duration.String(),
		Metrics:  domain.BuildSummary(metrics, duration),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) buildConfig(req testRequest) (config.Config, domain.Target, error) {
	cfg := s.cfg

	if req.Duration != "" {
		dur, err := time.ParseDuration(req.Duration)
		if err != nil {
			return config.Config{}, domain.Target{}, errors.New("invalid duration")
		}
		cfg.Duration = dur
	}

	if req.Concurrency != nil {
		cfg.Concurrency = *req.Concurrency
	}

	if req.Rate != nil {
		cfg.Rate = *req.Rate
	}

	if req.Mode != "" {
		cfg.Mode = req.Mode
	}

	if req.NoKeepAlive != nil {
		cfg.NoKeepAlive = *req.NoKeepAlive
	}

	targetURL := req.TargetURL
	if targetURL == "" {
		targetURL = cfg.TargetURL
	}
	if targetURL == "" {
		return config.Config{}, domain.Target{}, errors.New("target_url is required")
	}

	method := req.Method
	if method == "" {
		method = cfg.Method
	}
	if method == "" {
		method = http.MethodGet
	}

	headers := map[string]string{}
	for k, v := range cfg.Headers {
		headers[k] = v
	}
	for k, v := range req.Headers {
		headers[k] = v
	}

	target := domain.Target{
		URL:     targetURL,
		Method:  method,
		Headers: headers,
		Body:    []byte(req.Body),
	}

	return cfg, target, nil
}
