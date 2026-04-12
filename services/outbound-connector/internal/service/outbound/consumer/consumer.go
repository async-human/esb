package consumer

import (
	"context"

	"github.com/async-human/esb/platform/logger"
	"go.uber.org/zap"
)

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Outbound Connector consumer service")

	err := s.consumer.Consume(ctx, s.HandleMessage)
	if err != nil {
		logger.Error(ctx, "Consume error", zap.Error(err))
		return err
	}

	return nil
}