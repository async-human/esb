package consumer

import (
	"context"

	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/logger"
	"go.uber.org/zap"
)

func (s *service) HandleMessage(ctx context.Context, message kafka.Message) error {

	modelMsg, err := s.consumerDecoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderRecorded", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Message received",
		zap.String("topic", message.Topic),
		zap.Any("partition", message.Partition),
		zap.Any("offset", message.Offset),
		zap.String("message_id", modelMsg.Id.String()),
	)

	err = s.senderService.Send(ctx, modelMsg)
	if err != nil {
		logger.Warn(ctx, "failed to produce message", zap.Error(err), zap.String("message_id", modelMsg.Id.String()))
		return err
	}

	return nil
}