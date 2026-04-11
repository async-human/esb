package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/async-human/esb/inbound-connector/internal/model"
	"github.com/async-human/esb/platform/logger"
	"go.uber.org/zap"
)

func (p *service) ProduceMessageRecorded(ctx context.Context, message model.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	logger.Info(ctx, "📨 Sending message to Kafka",
		zap.String("message_id", message.Id.String()),
	)

	err = p.producer.Send(ctx, []byte(message.Id.String()), payload)
	if err != nil {
		logger.Error(ctx, "❌ Failed to send message to Kafka",
			zap.String("message_id", message.Id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	logger.Info(ctx, "✅ Message sent to Kafka",
		zap.String("message_id", message.Id.String()),
	)

	return nil
}
