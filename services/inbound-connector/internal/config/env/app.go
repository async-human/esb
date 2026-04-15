package env

import (
	"github.com/caarlos0/env/v11"
)

type appEnvConfig struct {
	AppName string `env:"APP_NAME" envDefault:"router-worker"`
	AppEnv  string `env:"APP_ENV" envDefault:"development"`
}

type appConfig struct {
	raw appEnvConfig
}

func NewAppConfig() (*appConfig, error) {
	var raw appEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &appConfig{raw: raw}, nil
}

func (cfg *appConfig) ServiceName() string {
	return cfg.raw.AppName
}

func (cfg *appConfig) Environment() string {
	return cfg.raw.AppEnv
}

func (cfg *appConfig) ServiceVersion() string {
	return "0.0.1"
}