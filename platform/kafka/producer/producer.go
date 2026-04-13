package producer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/async-human/esb/platform/kafka"
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
	middlewares  []kafka.ProducerMiddleware
}

func NewProducer(
	syncProducer sarama.SyncProducer,
	topics []string, logger Logger,
	middlewares ...kafka.ProducerMiddleware,
) *producer {
	return &producer{
		syncProducer: syncProducer,
		topics:       topics,
		logger:       logger,
		middlewares:  middlewares,
	}
}

func (p *producer) Send(ctx context.Context, msg kafka.Message) error {
	// Строим цепочку middleware вокруг базового handler
	handler := kafka.SendHandler(p.sendToTopics)
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		handler = p.middlewares[i](handler)
	}
	return handler(ctx, msg)
}

// sendToTopics — базовый handler без middleware
func (p *producer) sendToTopics(ctx context.Context, msg kafka.Message) error {
	saramaHeaders := toSaramaHeaders(msg.Headers)

	for _, topic := range p.topics {
		partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic:   topic,
			Key:     sarama.ByteEncoder(msg.Key),
			Value:   sarama.ByteEncoder(msg.Value),
			Headers: saramaHeaders,
		})
		if err != nil {
			p.logger.Error(ctx, "Failed to send message to Kafka", zap.Error(err)) // ← исправленный баг
			return err
		}
		p.logger.Info(ctx, "Message sent to Kafka",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
	}
	return nil
}

func toSaramaHeaders(headers map[string][]byte) []sarama.RecordHeader {
	result := make([]sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		result = append(result, sarama.RecordHeader{
			Key:   []byte(k),
			Value: v,
		})
	}
	return result
}
