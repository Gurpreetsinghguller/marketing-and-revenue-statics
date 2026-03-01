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
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type AuthConfig struct {
	SecretFile string `yaml:"secret_file"`
}

type RateLimitConfig struct {
	MaxRequests   int `yaml:"max_requests"`
	WindowSeconds int `yaml:"window_seconds"`
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{Port: "8080"},
		Log:    LogConfig{Level: "info"},
		Auth:   AuthConfig{SecretFile: "shared/secret"},
		RateLimit: RateLimitConfig{
			MaxRequests:   100,
			WindowSeconds: 60,
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

	return cfg, nil
}
