package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var(
	// Регистрация метра
	meter = otel.Meter("router-worker")

	// Starting app считает общее количество запусков
	AppStartsTotal, _ = meter.Int64Counter(
		"router_worker_starts_total",
		metric.WithDescription("Total number of starts Router Worker"),
	)

	// Starting app считает общее количество запусков
	AppEndTotal, _ = meter.Int64Counter(
		"router_worker_end_total",
		metric.WithDescription("Total number of end Router Worker"),
	)

)