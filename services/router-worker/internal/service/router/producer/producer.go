package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/logger"
	"github.com/async-human/esb/router-worker/internal/model"
	"go.uber.org/zap"
)

func (p *service) ProduceMessage(ctx context.Context, message model.Message) error {

	if message.Payload == nil {
		message.Payload = make(map[string]any)
	}
	message.Payload["test"] = "Hello from Router Worker!"

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.producer.Send(ctx, kafka.Message{
		Key:   []byte(message.Id.String()),
		Value: payload,
	})

	if err != nil {
		logger.Error(ctx, "❌ Failed to send message to Kafka",
			zap.String("message_id", message.Id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}
