package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type outboundProducerEnvConfig struct {
	Topics   []string `env:"KAFKA_OUTBOUND_TOPIC_NAME" envDefault:"esb.outbound.ready"`
}

type outboundProducerConfig struct {
	raw outboundProducerEnvConfig
}

func NewOutboundProducerConfig() (*outboundProducerConfig, error) {
	var raw outboundProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &outboundProducerConfig{raw: raw}, nil
}

func (cfg *outboundProducerConfig) Topics() []string {
	return cfg.raw.Topics
}

func (cfg *outboundProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
