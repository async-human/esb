package metrics

import (
	"errors"

	"go.opentelemetry.io/otel/metric"
)

// ── Subsystem structs ────────────────────────────────────

type App struct {
	StartsTotal metric.Int64Counter
	EndsTotal   metric.Int64Counter
	InFlight    metric.Int64UpDownCounter
}

type HTTP struct {
	RequestsTotal metric.Int64Counter
	Duration      metric.Float64Histogram
	RequestSize   metric.Int64Histogram
}

type KafkaProducer struct {
	MessagesTotal metric.Int64Counter
	Duration      metric.Float64Histogram
	Errors        metric.Int64Counter
}

type KafkaConsumer struct {
	MessagesTotal metric.Int64Counter
	ConsumerLag   metric.Int64Gauge
	FetchDuration metric.Float64Histogram
}

type Routing struct {
	DecisionsTotal metric.Int64Counter
	Duration       metric.Float64Histogram
	DLQTotal       metric.Int64Counter
	RetryTotal     metric.Int64Counter
}

type Delivery struct {
	AttemptsTotal metric.Int64Counter
	Duration      metric.Float64Histogram
	Errors        metric.Int64Counter
	RetryTotal    metric.Int64Counter
}

// ── Constructors ─────────────────────────────────────────
// Принимают meter снаружи. Возвращают ошибку.

func NewApp(m metric.Meter) (App, error) {
	starts, e1 := m.Int64Counter("app.starts",
		metric.WithDescription("Total application starts"),
		metric.WithUnit("{start}"))
	ends, e2 := m.Int64Counter("app.ends",
		metric.WithDescription("Total application ends"),
		metric.WithUnit("{end}"))
	flight, e3 := m.Int64UpDownCounter("app.messages.in_flight",
		metric.WithDescription("Messages currently in processing pipeline"),
		metric.WithUnit("{message}"))

	return App{starts, ends, flight}, errors.Join(e1, e2, e3)
}

func NewHTTP(m metric.Meter) (HTTP, error) {
	requests, e1 := m.Int64Counter("http.server.request.count",
		metric.WithDescription("Total HTTP requests received"),
		metric.WithUnit("{request}"))
	duration, e2 := m.Float64Histogram("http.server.request.duration",
		metric.WithDescription("HTTP server request duration"),
		metric.WithUnit("s"))
	size, e3 := m.Int64Histogram("http.server.request.body.size",
		metric.WithDescription("HTTP request body size"),
		metric.WithUnit("By"))

	return HTTP{requests, duration, size}, errors.Join(e1, e2, e3)
}

func NewKafkaProducer(m metric.Meter) (KafkaProducer, error) {
	msgs, e1 := m.Int64Counter("kafka.producer.messages",
		metric.WithDescription("Total messages produced"),
		metric.WithUnit("{message}"))
	dur, e2 := m.Float64Histogram("kafka.producer.duration",
		metric.WithDescription("Kafka produce latency"),
		metric.WithUnit("s"))
	errs, e3 := m.Int64Counter("kafka.producer.errors",
		metric.WithDescription("Kafka produce errors"),
		metric.WithUnit("{error}"))

	return KafkaProducer{msgs, dur, errs}, errors.Join(e1, e2, e3)
}

func NewKafkaConsumer(m metric.Meter) (KafkaConsumer, error) {
	msgs, e1 := m.Int64Counter("kafka.consumer.messages",
		metric.WithDescription("Total messages consumed"),
		metric.WithUnit("{message}"))
	lag, e2 := m.Int64Gauge("kafka.consumer.lag",
		metric.WithDescription("Consumer lag — key SLO metric"),
		metric.WithUnit("{message}"))
	dur, e3 := m.Float64Histogram("kafka.consumer.fetch.duration",
		metric.WithDescription("Kafka fetch batch duration"),
		metric.WithUnit("s"))

	return KafkaConsumer{msgs, lag, dur}, errors.Join(e1, e2, e3)
}

func NewRouting(m metric.Meter) (Routing, error) {
	dec, e1 := m.Int64Counter("routing.decisions",
		metric.WithDescription("Routing decisions made"),
		metric.WithUnit("{decision}"))
	dur, e2 := m.Float64Histogram("routing.decision.duration",
		metric.WithDescription("Time to make a routing decision"),
		metric.WithUnit("s"))
	dlq, e3 := m.Int64Counter("routing.dlq.messages",
		metric.WithDescription("Messages sent to Dead Letter Queue"),
		metric.WithUnit("{message}"))
	ret, e4 := m.Int64Counter("routing.retries",
		metric.WithDescription("Routing retry attempts"),
		metric.WithUnit("{attempt}"))

	return Routing{dec, dur, dlq, ret}, errors.Join(e1, e2, e3, e4)
}

func NewDelivery(m metric.Meter) (Delivery, error) {
	att, e1 := m.Int64Counter("delivery.attempts",
		metric.WithDescription("Total outbound delivery attempts"),
		metric.WithUnit("{attempt}"))
	dur, e2 := m.Float64Histogram("delivery.duration",
		metric.WithDescription("Outbound delivery latency"),
		metric.WithUnit("s"))
	errs, e3 := m.Int64Counter("delivery.errors",
		metric.WithDescription("Outbound delivery errors"),
		metric.WithUnit("{error}"))
	ret, e4 := m.Int64Counter("delivery.retries",
		metric.WithDescription("Outbound delivery retry attempts"),
		metric.WithUnit("{attempt}"))

	return Delivery{att, dur, errs, ret}, errors.Join(e1, e2, e3, e4)
}