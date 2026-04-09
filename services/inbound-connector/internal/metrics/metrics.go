package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var(
	// Регистрация метра
	meter = otel.Meter("inbound-connector")

	// Starting app считает общее количество запусков
	AppStartsTotal, _ = meter.Int64Counter(
		"inbound_connector_starts_total",
		metric.WithDescription("Total number of starts Inbound Connector"),
	)

	// Starting app считает общее количество запусков
	AppEndTotal, _ = meter.Int64Counter(
		"inbound_connector_end_total",
		metric.WithDescription("Total number of end Inbound Connector"),
	)

)