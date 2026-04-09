package config

import (
	"os"

	"github.com/async-human/esb/outbound-connector/internal/config/env"
	"github.com/joho/godotenv"
)

var commonAppConfig *config

type config struct {
	Logger       LoggerConfig
	App          AppConfig
	Otel         OtelConfig
	MetricConfig MetricServerConfig
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

	metricCfg, err := env.NewMetricConfig()
	if err != nil {
		return err
	}

	commonAppConfig = &config{
		Logger:       loggerCfg,
		App:          appCfg,
		Otel:         otelCfg,
		MetricConfig: metricCfg,
	}

	return nil

}

func CommonAppConfig() *config {
	return commonAppConfig
}
