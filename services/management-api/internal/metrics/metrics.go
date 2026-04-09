package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var(
	// Регистрация метра
	meter = otel.Meter("management-api")

	// Starting app считает общее количество запусков
	AppStartsTotal, _ = meter.Int64Counter(
		"management_api_starts_total",
		metric.WithDescription("Total number of starts Management API"),
	)

	// Starting app считает общее количество запусков
	AppEndTotal, _ = meter.Int64Counter(
		"management_api_end_total",
		metric.WithDescription("Total number of end Management API"),
	)

)