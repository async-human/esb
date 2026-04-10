package env

import (
	"github.com/caarlos0/env/v11"
)

type KafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" envDefault:"kafka:9092"`
}

type KafkaConfig struct {
	raw KafkaEnvConfig
}

func NewKafkaConfig() (*KafkaConfig, error) {
	var raw KafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &KafkaConfig{raw: raw}, nil
}

func (cfg *KafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}
