package consumer

import (
	"context"

	"github.com/async-human/esb/platform/logger"
	"github.com/async-human/esb/router-worker/internal/metrics"
	"go.uber.org/zap"
)

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Router Worker consumer service")

	metrics.ServiceMetrics.App.StartsTotal.Add(ctx, 1)
	defer metrics.ServiceMetrics.App.EndsTotal.Add(ctx, 1)

	err := s.consumer.Consume(ctx, s.HandleMessage)
	if err != nil {
		logger.Error(ctx, "Consume error", zap.Error(err),
			zap.Error(err),
		)
		return err
	}

	return nil
}
