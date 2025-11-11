package infrastructure

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/net/proxy"

	"go-load-lab/config"
	"go-load-lab/domain"
)

type Executor struct {
	clients []*http.Client
	index   uint32 // for round-robin
}

func defaultTransport(cfg config.Config) *http.Transport {
	return &http.Transport{
		MaxIdleConns:       100_000,
		MaxConnsPerHost:    100_000,
		DisableKeepAlives:  cfg.NoKeepAlive,
		DisableCompression: false,
	}
}

func NewExecutor(cfg config.Config) *Executor {
	e := &Executor{}

	transports := []*http.Transport{}

	// no-proxy client всегда
	transports = append(transports, defaultTransport(cfg))

	for _, p := range cfg.Proxies {
		u, err := url.Parse(p)
		if err != nil {
			continue
		}

		tr := defaultTransport(cfg)

		switch u.Scheme {
		case "http", "https":
			tr.Proxy = http.ProxyURL(u)
		case "socks5", "socks5h":
			dialer, err := proxy.FromURL(u, proxy.Direct)
			if err != nil {
				continue
			}
			if cd, ok := dialer.(proxy.ContextDialer); ok {
				tr.DialContext = cd.DialContext
			} else {
				tr.Dial = dialer.Dial
			}
		default:
			continue
		}

		transports = append(transports, tr)
	}

	for _, tr := range transports {
		e.clients = append(e.clients, &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		})
	}

	if len(e.clients) == 0 {
		e.clients = append(e.clients, &http.Client{Transport: defaultTransport(cfg)})
	}

	return e
}

func (e *Executor) getClient() *http.Client {
	i := atomic.AddUint32(&e.index, 1) - 1
	return e.clients[i%uint32(len(e.clients))]
}

func (e *Executor) Do(target domain.Target) domain.Result {
	client := e.getClient()

	req, err := http.NewRequest(target.Method, target.URL, nil)
	if err != nil {
		return domain.Result{Err: err}
	}

	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}

	if len(target.Body) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(target.Body))
		req.ContentLength = int64(len(target.Body))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(target.Body)), nil
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return domain.Result{Latency: latency, Err: err}
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return domain.Result{Latency: latency, Status: resp.StatusCode}
}
