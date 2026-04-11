package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OtelConfig interface {
	OTelCollectorEndpoint() string
	OTelEnabled() bool
}

type AppConfig interface {
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type MetricServerConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

type RestConfig interface {
	Host() string
	Port() string
	Address() string
}

type KafkaConfig interface {
	Brokers() []string
}

type InboundProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}
