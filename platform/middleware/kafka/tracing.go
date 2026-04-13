package kafka

import (
	"context"

	platformKafka "github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/kafka/consumer"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "kafka"

// TracingConsumer — извлекает TraceContext из headers входящего сообщения.
// Создаёт дочерний span, помещает его в context.
func TracingConsumer() consumer.Middleware {
	tracer := otel.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return func(next consumer.MessageHandler) consumer.MessageHandler {
		return func(ctx context.Context, msg platformKafka.Message) error {
			// Extract: восстанавливаем родительский контекст из headers
			carrier := kafkaHeaderCarrier(msg.Headers)
			ctx = propagator.Extract(ctx, carrier)

			// Создаём дочерний span для обработки этого сообщения
			ctx, span := tracer.Start(ctx, "kafka.consume",
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					semconv.MessagingSystemKey.String("kafka"),
					semconv.MessagingDestinationNameKey.String(msg.Topic),
					semconv.MessagingKafkaMessageOffsetKey.Int(int(msg.Offset)),
				),
			)
			defer span.End()

			err := next(ctx, msg)
			if err != nil {
				span.RecordError(err)
			}
			return err
		}
	}
}

// TracingProducer — инжектирует TraceContext из context в headers исходящего сообщения.
func TracingProducer() platformKafka.ProducerMiddleware {
	tracer := otel.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return func(next platformKafka.SendHandler) platformKafka.SendHandler {
		return func(ctx context.Context, msg platformKafka.Message) error {
			// Создаём span для отправки
			ctx, span := tracer.Start(ctx, "kafka.produce",
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					semconv.MessagingSystemKey.String("kafka"),
				),
			)
			defer span.End()

			// Inject: записываем traceparent + tracestate в headers сообщения
			if msg.Headers == nil {
				msg.Headers = make(map[string][]byte)
			}
			carrier := kafkaHeaderCarrier(msg.Headers)
			propagator.Inject(ctx, carrier)
			// msg.Headers теперь содержит, например:
			// "traceparent" -> "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
			// "tracestate"  -> ""

			err := next(ctx, msg)
			if err != nil {
				span.RecordError(err)
			}
			return err
		}
	}
}