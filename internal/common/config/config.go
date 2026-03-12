package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "config/config.yml"

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Log       LogConfig       `yaml:"log"`
	Auth      AuthConfig      `yaml:"auth"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Broker    BrokerConfig    `yaml:"broker"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type AuthConfig struct {
	SecretFile string `yaml:"secret_file"`
	TokenTTL   string `yaml:"token_ttl"`
}

type RateLimitConfig struct {
	MaxRequests   int `yaml:"max_requests"`
	WindowSeconds int `yaml:"window_seconds"`
}

type BrokerConfig struct {
	Type     string `yaml:"type"`      // "mqtt", "kafka", etc.
	URL      string `yaml:"url"`       // Broker connection URL
	ClientID string `yaml:"client_id"` // Client ID for broker
	Topic    string `yaml:"topic"`     // Topic to subscribe to
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{Port: "8080"},
		Log:    LogConfig{Level: "info"},
		Auth: AuthConfig{
			SecretFile: "shared/secret",
			TokenTTL:   "24h"},
		RateLimit: RateLimitConfig{
			MaxRequests:   100,
			WindowSeconds: 60,
		},
		Broker: BrokerConfig{
			Type:     "mqtt",
			URL:      "tcp://localhost:1883",
			ClientID: "marketing-app",
			Topic:    "events",
		},
	}
}

func Load(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultConfigPath
	}

	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if len(data) == 0 {
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if strings.TrimSpace(cfg.Server.Port) == "" {
		cfg.Server.Port = "8080"
	}
	if strings.TrimSpace(cfg.Log.Level) == "" {
		cfg.Log.Level = "info"
	}
	if strings.TrimSpace(cfg.Auth.SecretFile) == "" {
		cfg.Auth.SecretFile = "shared/secret"
	}
	if cfg.RateLimit.MaxRequests <= 0 {
		cfg.RateLimit.MaxRequests = 100
	}
	if cfg.RateLimit.WindowSeconds <= 0 {
		cfg.RateLimit.WindowSeconds = 60
	}
	if strings.TrimSpace(cfg.Broker.Type) == "" {
		cfg.Broker.Type = "mqtt"
	}
	if strings.TrimSpace(cfg.Broker.URL) == "" {
		cfg.Broker.URL = "tcp://localhost:1883"
	}
	if strings.TrimSpace(cfg.Broker.ClientID) == "" {
		cfg.Broker.ClientID = "marketing-app"
	}
	if strings.TrimSpace(cfg.Broker.Topic) == "" {
		cfg.Broker.Topic = "events"
	}

	return cfg, nil
}
