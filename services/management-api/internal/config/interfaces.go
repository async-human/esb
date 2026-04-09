package config

import "time"

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