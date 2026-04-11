package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/async-human/esb/platform/closer"
	platformKafka "github.com/async-human/esb/platform/kafka"
	platformKafkaConsumer "github.com/async-human/esb/platform/kafka/consumer"
	platformKafkaProducer "github.com/async-human/esb/platform/kafka/producer"
	"github.com/async-human/esb/platform/logger"
	"github.com/async-human/esb/router-worker/internal/config"
	kafkaConverter "github.com/async-human/esb/router-worker/internal/converter/kafka"
	"github.com/async-human/esb/router-worker/internal/converter/kafka/decode"
	"github.com/async-human/esb/router-worker/internal/service"
	"github.com/async-human/esb/router-worker/internal/service/router/consumer"
	producerService "github.com/async-human/esb/router-worker/internal/service/router/producer"
)

type diContainer struct {
	consumerService service.ConsumerService
	producerService service.ProducerService

	consumer        platformKafka.Consumer
	consumerGroup   sarama.ConsumerGroup
	consumerDecoder kafkaConverter.MessageDecoder

	syncProducer     sarama.SyncProducer
	platformProducer platformKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerService() service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = consumer.NewService(
			d.Consumer(),
			d.ProducerService(),
			d.ConsumerDecoder(),
		)
	}
	return d.consumerService
}

func (d *diContainer) ProducerService() service.ProducerService {
	if d.producerService == nil {
		d.producerService = producerService.NewService(
			d.PlatformProducer(),
		)
	}
	return d.producerService
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
		)
	}
	return d.consumer
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

func (d *diContainer) ConsumerDecoder() kafkaConverter.MessageDecoder {
	if d.consumerDecoder == nil {
		d.consumerDecoder = decode.NewDecoder()
	}
	return d.consumerDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		syncProducer, err := sarama.NewSyncProducer(
			config.CommonAppConfig().KafkaConfig.Brokers(),
			config.CommonAppConfig().OutboundProducerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return syncProducer.Close()
		})
		d.syncProducer = syncProducer
	}
	return d.syncProducer
}

func (d *diContainer) PlatformProducer() platformKafka.Producer {
	if d.platformProducer == nil {
		d.platformProducer = platformKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.CommonAppConfig().OutboundProducerConfig.Topics(),
			logger.Logger(),
		)
	}
	return d.platformProducer
}
