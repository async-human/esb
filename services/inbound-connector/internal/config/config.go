package config

import (
	"os"

	"github.com/async-human/esb/inbound-connector/internal/config/env"
	"github.com/joho/godotenv"
)

var commonAppConfig *config

type config struct {
	Logger LoggerConfig
	App    AppConfig
	Otel   OtelConfig
}

func Load(path ...string) error {

	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appCfg, err := env.NewAppConfig()
	if err != nil {
		return err
	}

	otelCfg, err := env.NewOtelConfig()
	if err != nil {
		return err
	}

	commonAppConfig = &config{
		Logger: loggerCfg,
		App: appCfg,
		Otel: otelCfg,
	}

	return nil

}

func CommonAppConfig() *config {
	return commonAppConfig
}
