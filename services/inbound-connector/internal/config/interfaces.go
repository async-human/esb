package config

import (
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