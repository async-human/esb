package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type inboundConsumerEnvConfig struct {
	Topic  string `env:"KAFKA_INBOUND_TOPIC_NAME" envDefault:"esb.inbound.raw"`
	GroupID string   `env:"KAFKA_INBOUND_CONSUMER_GROUP_ID" envDefault:"esb-inbound-group-raw"`
}

type inboundConsumerConfig struct {
	raw inboundConsumerEnvConfig
}

func NewInboundConsumerConfig() (*inboundConsumerConfig, error) {
	var raw inboundConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inboundConsumerConfig{raw: raw}, nil
}

func (cfg *inboundConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *inboundConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *inboundConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
