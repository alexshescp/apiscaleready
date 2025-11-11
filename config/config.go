package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	TargetURL      string
	Method         string
	Concurrency    int
	Duration       time.Duration
	Rate           int // req/s, 0 = unlimited
	Mode           string
	NoKeepAlive    bool
	ProxyFile      string
	JsonBody       string
	HeadersJSON    string
	ContentType    string
	Headers        map[string]string
	Proxies        []string
}

func Load() Config {
	_ = godotenv.Load()

	viper.SetEnvPrefix("GLAB")
	viper.AutomaticEnv()

	viper.SetDefault("METHOD", "GET")
	viper.SetDefault("CONCURRENCY", 100)
	viper.SetDefault("DURATION", "10s")
	viper.SetDefault("RATE", 0)
	viper.SetDefault("MODE", "wg")
	viper.SetDefault("NO_KEEP_ALIVE", false)
	viper.SetDefault("CONTENT_TYPE", "application/json")

	dur, err := time.ParseDuration(viper.GetString("DURATION"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		TargetURL:   viper.GetString("TARGET_URL"),
		Method:      viper.GetString("METHOD"),
		Concurrency: viper.GetInt("CONCURRENCY"),
		Duration:    dur,
		Rate:        viper.GetInt("RATE"),
		Mode:        viper.GetString("MODE"),
		NoKeepAlive: viper.GetBool("NO_KEEP_ALIVE"),
		ProxyFile:   viper.GetString("PROXY_FILE"),
		JsonBody:    viper.GetString("JSON_BODY"),
		HeadersJSON: viper.GetString("HEADERS_JSON"),
		ContentType: viper.GetString("CONTENT_TYPE"),
	}

	if cfg.JsonBody != "" && cfg.ContentType == "" {
		cfg.ContentType = "application/json"
	}

	if cfg.HeadersJSON != "" {
		cfg.Headers = make(map[string]string)
		_ = json.Unmarshal([]byte(cfg.HeadersJSON), &cfg.Headers)
	}

	if cfg.ContentType != "" {
		if cfg.Headers == nil {
			cfg.Headers = make(map[string]string)
		}
		cfg.Headers["Content-Type"] = cfg.ContentType
	}

	if cfg.ProxyFile != "" {
		data, err := os.ReadFile(cfg.ProxyFile)
		if err != nil {
			log.Fatal(err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				cfg.Proxies = append(cfg.Proxies, line)
			}
		}
	}

	return cfg
}