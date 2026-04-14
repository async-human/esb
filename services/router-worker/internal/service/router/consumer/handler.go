package consumer

import (
	"context"
	"time"

	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/logger"
	"github.com/async-human/esb/router-worker/internal/metrics"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

func (s *service) HandleMessage(ctx context.Context, message kafka.Message) error {
	start := time.Now()

	// Сколько сообщений обрабатывается прямо сейчас
	metrics.ServiceMetrics.App.InFlight.Add(ctx, 1)
	defer metrics.ServiceMetrics.App.InFlight.Add(ctx, -1)

	// Сообщение получено из Kafka
	metrics.ServiceMetrics.Consumer.MessagesTotal.Add(ctx, 1,
		otelmetric.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.topic", message.Topic),
			attribute.Int("messaging.partition", int(message.Partition)),
		),
	)

	modelMsg, err := s.consumerDecoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode message", zap.Error(err))

		// Decode failure — это routing error: сообщение пришло, но не может быть обработано
		metrics.ServiceMetrics.Routing.DecisionsTotal.Add(ctx, 1,
			otelmetric.WithAttributes(
				attribute.String("routing.result", "decode_error"),
				attribute.String("messaging.topic", message.Topic),
			),
		)
		metrics.ServiceMetrics.Routing.DLQTotal.Add(ctx, 1,
			otelmetric.WithAttributes(
				attribute.String("routing.dlq_reason", "decode_error"),
			),
		)
		return err
	}

	logger.Info(ctx, "Message received",
		zap.String("topic", message.Topic),
		zap.Any("partition", message.Partition),
		zap.Any("offset", message.Offset),
		zap.String("message_id", modelMsg.Id.String()),
	)

	err = s.producerService.ProduceMessage(ctx, modelMsg)

	// Определяем результат маршрутизации до записи метрик
	routingResult := "success"
	if err != nil {
		routingResult = "produce_error"
	}

	routingAttrs := otelmetric.WithAttributes(
		attribute.String("routing.result", routingResult),
		attribute.String("messaging.topic", message.Topic),
	)

	// Exemplar прикрепится из span'а TracingConsumer middleware
	metrics.ServiceMetrics.Routing.Duration.Record(ctx, time.Since(start).Seconds(), routingAttrs)
	metrics.ServiceMetrics.Routing.DecisionsTotal.Add(ctx, 1, routingAttrs)

	if err != nil {
		logger.Warn(ctx, "Failed to produce message",
			zap.Error(err),
			zap.String("message_id", modelMsg.Id.String()),
		)
		metrics.ServiceMetrics.Routing.DLQTotal.Add(ctx, 1,
			otelmetric.WithAttributes(
				attribute.String("routing.dlq_reason", "produce_error"),
				attribute.String("messaging.topic", message.Topic),
			),
		)
		return err
	}

	return nil
}