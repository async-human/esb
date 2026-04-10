package env

import (
	"github.com/caarlos0/env/v11"
)

type restEnvConfig struct {
	Host    string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port 	string  `env:"SERVER_PORT" envDefault:"8080"`
}

type restConfig struct {
	raw restEnvConfig
}

func NewRestConfig() (*restConfig, error) {
	var raw restEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &restConfig{raw: raw}, nil
}

func (cfg *restConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *restConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *restConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}