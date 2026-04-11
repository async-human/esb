package producer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/async-human/esb/platform/logger"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producer struct {
	syncProducer sarama.SyncProducer
	topics       []string
	logger       Logger
}

func NewProducer(syncProducer sarama.SyncProducer, topics []string, logger Logger) *producer {
	return &producer{
		syncProducer: syncProducer,
		topics:       topics,
		logger:       logger,
	}
}

func (p *producer) Send(ctx context.Context, key, value []byte) error {
	
	for _, topic := range p.topics {
		partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(value),
		})
		if err != nil {
			logger.Error(ctx, "Failed send message to Kafka", zap.Error(err))
			return err
		}
		p.logger.Info(ctx, "Message sent to Kafka",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
			zap.String("key", string(key)),
			zap.String("value", string(value)),
		)
	}

	return nil
}
