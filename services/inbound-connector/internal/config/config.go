package config

import (
	"os"

	"github.com/async-human/esb/inbound-connector/internal/config/env"
	"github.com/joho/godotenv"
)

var commonAppConfig *config

type config struct {
	Logger       LoggerConfig
	App          AppConfig
	Otel         OtelConfig
	MetricConfig MetricServerConfig
	Rest         RestConfig
	Kafka        KafkaConfig
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

	restCfg, err := env.NewRestConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	commonAppConfig = &config{
		Logger:       loggerCfg,
		App:          appCfg,
		Otel:         otelCfg,
		MetricConfig: metricCfg,
		Rest:         restCfg,
		Kafka:        kafkaCfg,
	}

	return nil

}

func CommonAppConfig() *config {
	return commonAppConfig
}
