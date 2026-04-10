package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/async-human/esb/platform/kafka"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      Logger
	middlewares []Middleware
}

// NewConsumer — создаёт новый consumer.
func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

// Consume запускает консьюмер для списка топиков.
func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	wrappedHandler := func(ctx context.Context, msg kafka.Message) error {
		kafkaMsg := kafka.Message{
			Headers:        msg.Headers,
			Timestamp:      msg.Timestamp,
			BlockTimestamp: msg.BlockTimestamp,
			Key:            msg.Key,
			Value:          msg.Value,
			Topic:          msg.Topic,
			Partition:      msg.Partition,
			Offset:         msg.Offset,
		}

		return handler(ctx, kafkaMsg)
	}

	newGroupHandler := NewGroupHandler(wrappedHandler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Kafka consume error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}