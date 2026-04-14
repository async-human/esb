package metrics

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
)

const (
	defaultTimeout = 5 * time.Second
)

var (
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *metric.MeterProvider
)

type Config interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

// InitProvider инициализирует глобальный провайдер метрик OpenTelemetry
func InitProvider(ctx context.Context, cfg Config) error {
	var err error

	// Создаем экспортер для отправки метрик в OTLP коллектор
	exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(defaultTimeout),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	// Создаем провайдер метрик
	meterProvider = metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(cfg.CollectorInterval()),
			),
		),
		metric.WithExemplarFilter(exemplar.AlwaysOnFilter),
		metric.WithView(views()...),
	)

	// Устанавливаем глобальный провайдер метрик
	otel.SetMeterProvider(meterProvider)

	return nil
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

// ForceFlush принудительно отправляет все накопленные метрики в OTel Collector
func ForceFlush(ctx context.Context) error {
	if meterProvider == nil {
		return nil
	}
	return errors.Wrap(meterProvider.ForceFlush(ctx), "failed to force flush metrics")
}

// Shutdown закрывает провайдер метрик и экспортер
func Shutdown(ctx context.Context) error {
	if meterProvider == nil && exporter == nil {
		return nil
	}

	var err error

	if meterProvider != nil {
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
	}

	if exporter != nil {
		err = exporter.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown exporter")
		}
	}

	return nil
}

var histogramBoundaries = map[string][]float64{
	"http.server.request.duration":  {.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	"kafka.producer.duration":       {.001, .005, .01, .025, .05, .1, .25, .5, 1},
	"kafka.consumer.fetch.duration": {.001, .005, .01, .025, .05, .1, .25, .5, 1},
	"routing.decision.duration":     {.0001, .001, .005, .01, .025, .05, .1},
	"delivery.duration":             {.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
}

func views() []metric.View {
	views := make([]metric.View, 0, len(histogramBoundaries))
	for name, bounds := range histogramBoundaries {
		views = append(views, metric.NewView(
			metric.Instrument{Name: name},
			metric.Stream{
				Aggregation: metric.AggregationExplicitBucketHistogram{
					Boundaries: bounds,
				},
			},
		))
	}
	return views
}
