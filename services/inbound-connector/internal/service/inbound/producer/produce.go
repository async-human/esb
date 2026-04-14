package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/async-human/esb/inbound-connector/internal/metrics"
	"github.com/async-human/esb/inbound-connector/internal/model"
	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/logger"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

func (p *service) ProduceMessage(ctx context.Context, message model.Message) error {
	start := time.Now()

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	sendErr := p.producer.Send(ctx, kafka.Message{
		Key:   []byte(message.Id.String()),
		Value: payload,
	})

	// Общий набор атрибутов для всех метрик этого вызова
	attrs := []attribute.KeyValue{
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.operation", "publish"),
	}

	status := "success"
	if sendErr != nil {
		status = "error"
	}

	attrsWithStatus := otelmetric.WithAttributes(
		append(attrs, attribute.String("messaging.status", status))...,
	)

	// Exemplar прикрепится автоматически из span'а TracingProducer middleware
	metrics.ServiceMetrics.Producer.Duration.Record(ctx, time.Since(start).Seconds(), attrsWithStatus)
	metrics.ServiceMetrics.Producer.MessagesTotal.Add(ctx, 1, attrsWithStatus)

	if sendErr != nil {
		metrics.ServiceMetrics.Producer.Errors.Add(ctx, 1,
			otelmetric.WithAttributes(
				append(attrs, attribute.String("messaging.error_type", sendErr.Error()))...,
			),
		)

		logger.Error(ctx, "❌ Failed to send message to Kafka",
			zap.String("message_id", message.Id.String()),
			zap.Error(sendErr),
		)
		return fmt.Errorf("failed to send message to kafka: %w", sendErr)
	}

	return nil
}