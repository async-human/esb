package kafka

import (
	"context"

	"github.com/async-human/esb/platform/kafka"
	"github.com/async-human/esb/platform/kafka/consumer"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
}

func LoggingConsumer(logger Logger) consumer.Middleware {
	return func(next consumer.MessageHandler) consumer.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			logger.Info(ctx, "Kafka msg received", zap.String("topic", msg.Topic))
			return next(ctx, msg)
		}
	}
}

func LoggingProducer(logger Logger) kafka.ProducerMiddleware {
    return func(next kafka.SendHandler) kafka.SendHandler {
        return func(ctx context.Context, msg kafka.Message) error {
            logger.Info(ctx, "Kafka msg sending", zap.String("topic", msg.Topic))
            return next(ctx, msg)
        }
    }
}