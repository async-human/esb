package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/async-human/esb/outbound-connector/internal/config"
	kafkaConverter "github.com/async-human/esb/outbound-connector/internal/converter/kafka"
	platformKafkaMiddleware "github.com/async-human/esb/platform/middleware/kafka"
	"github.com/async-human/esb/outbound-connector/internal/converter/kafka/decode"
	"github.com/async-human/esb/outbound-connector/internal/service"
	"github.com/async-human/esb/outbound-connector/internal/service/outbound/consumer"
	httpSender "github.com/async-human/esb/outbound-connector/internal/service/outbound/http"
	"github.com/async-human/esb/platform/closer"
	platformKafka "github.com/async-human/esb/platform/kafka"
	platformKafkaConsumer "github.com/async-human/esb/platform/kafka/consumer"
	"github.com/async-human/esb/platform/logger"
)

type diContainer struct {
	consumerService service.ConsumerService
	senderService   service.SenderService

	consumer        platformKafka.Consumer
	consumerGroup   sarama.ConsumerGroup
	consumerDecoder kafkaConverter.MessageDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerService() service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = consumer.NewService(
			d.Consumer(),
			config.CommonAppConfig().InboundConsumerConfig.Endpoints(),
			d.SenderService(),
			d.ConsumerDecoder(),
		)
	}
	return d.consumerService
}

func (d *diContainer) Consumer() platformKafka.Consumer {
	if d.consumer == nil {
		topics := []string{
			config.CommonAppConfig().InboundConsumerConfig.Topic(),
		}
		d.consumer = platformKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			topics,
			logger.Logger(),
			platformKafkaMiddleware.TracingConsumer(), // ← сначала tracing
			platformKafkaMiddleware.LoggingConsumer(logger.Logger()),
		)
	}
	return d.consumer
}

func (d *diContainer) ConsumerDecoder() kafkaConverter.MessageDecoder {
	if d.consumerDecoder == nil {
		d.consumerDecoder = decode.NewDecoder()
	}
	return d.consumerDecoder
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.CommonAppConfig().KafkaConfig.Brokers(),
			config.CommonAppConfig().InboundConsumerConfig.GroupID(),
			config.CommonAppConfig().InboundConsumerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})
		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) SenderService() service.SenderService {
	if d.senderService == nil {
		d.senderService = httpSender.NewService()
	}
	return d.senderService
}
