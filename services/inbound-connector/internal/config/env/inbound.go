package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type inboundProducerEnvConfig struct {
	Topic   string `env:"KAFKA_INBOUND_TOPIC_NAME" envDefault:"esb.inbound.raw"`
}

type inboundProducerConfig struct {
	raw inboundProducerEnvConfig
}

func NewInboundProducerConfig() (*inboundProducerConfig, error) {
	var raw inboundProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inboundProducerConfig{raw: raw}, nil
}

func (cfg *inboundProducerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *inboundProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	return config
}
