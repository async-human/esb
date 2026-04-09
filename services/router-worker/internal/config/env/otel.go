package env

import (
	"github.com/caarlos0/env/v11"
)

type otelEnvConfig struct {
	Endpoint    string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
	OTelEnabled bool   `env:"OTEL_ENABLED" envDefault:"true"`
}

type otelConfig struct {
	raw otelEnvConfig
}

func NewOtelConfig() (*otelConfig, error) {
	var raw otelEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &otelConfig{raw: raw}, nil
}

func (cfg *otelConfig) OTelCollectorEndpoint() string {
	return cfg.raw.Endpoint
}

func (cfg *otelConfig) OTelEnabled() bool {
	return cfg.raw.OTelEnabled
}
