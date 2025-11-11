package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-load-lab/config"
)

func TestRunTestEndpoint(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer target.Close()

	cfg := config.Config{
		TargetURL:   target.URL,
		Method:      http.MethodGet,
		Concurrency: 1,
		Duration:    20 * time.Millisecond,
		Mode:        "wg",
	}

	server := New(cfg)
	api := httptest.NewServer(server.Handler())
	defer api.Close()

	body := map[string]string{"duration": "20ms"}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(api.URL+"/tests", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	var out struct {
		Metrics struct {
			Total   uint64 `json:"total"`
			Success uint64 `json:"success"`
		} `json:"metrics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if out.Metrics.Total == 0 {
		t.Fatalf("expected total requests to be > 0")
	}
	if out.Metrics.Success == 0 {
		t.Fatalf("expected successful requests to be > 0")
	}
}
