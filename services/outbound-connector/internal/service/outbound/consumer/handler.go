package consumer

import (
	"context"
	"time"

	"github.com/async-human/esb/outbound-connector/internal/metrics"
	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/logger"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

func (s *service) HandleMessage(ctx context.Context, message kafka.Message) error {
	start := time.Now()

	// Сколько сообщений в обработке прямо сейчас
	metrics.ServiceMetrics.App.InFlight.Add(ctx, 1)
	defer metrics.ServiceMetrics.App.InFlight.Add(ctx, -1)

	// Факт получения из Kafka — до decode, чтобы считать даже битые сообщения
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

		// Decode failure считается как delivery.attempts с результатом decode_error —
		// сообщение было получено, но доставить его невозможно
		metrics.ServiceMetrics.Delivery.AttemptsTotal.Add(ctx, 1,
			otelmetric.WithAttributes(
				attribute.String("delivery.result", "decode_error"),
				attribute.String("messaging.topic", message.Topic),
			),
		)
		metrics.ServiceMetrics.Delivery.Errors.Add(ctx, 1,
			otelmetric.WithAttributes(
				attribute.String("delivery.error_type", "decode_error"),
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

	// send.go запишет delivery.duration и delivery.attempts сам —
	// здесь только фиксируем общее время обработки сообщения handler'ом
	err = s.senderService.Send(ctx, modelMsg)

	result := "success"
	if err != nil {
		result = "send_error"
	}

	// Exemplar прикрепится из span'а TracingConsumer middleware
	metrics.ServiceMetrics.Consumer.FetchDuration.Record(ctx, time.Since(start).Seconds(),
		otelmetric.WithAttributes(
			attribute.String("messaging.topic", message.Topic),
			attribute.String("messaging.result", result),
		),
	)

	if err != nil {
		logger.Warn(ctx, "Failed to send message",
			zap.Error(err),
			zap.String("message_id", modelMsg.Id.String()),
		)
		return err
	}

	return nil
}