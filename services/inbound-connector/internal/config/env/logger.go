package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level    string `env:"LOGGER_LEVEL" envDefault:"info"`
	AsJson   bool   `env:"LOGGER_AS_JSON" envDefault:"true"`
	Endpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}
