package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	api "github.com/async-human/esb/inbound-connector/internal/api/v1"
	"github.com/async-human/esb/inbound-connector/internal/config"
	servicedDef "github.com/async-human/esb/inbound-connector/internal/service"
	inbound "github.com/async-human/esb/inbound-connector/internal/service/inbound"
	inboundproducer "github.com/async-human/esb/inbound-connector/internal/service/inbound/producer"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
	"github.com/async-human/esb/platform/closer"
	"github.com/async-human/esb/platform/kafka"
	kafkaproducer "github.com/async-human/esb/platform/kafka/producer"
	platformKafkaMiddleware "github.com/async-human/esb/platform/middleware/kafka"
	"github.com/async-human/esb/platform/logger"
)

type diContainer struct {
	messageProducerService servicedDef.ProducerService
	messageHandler         icv1.ServerInterface
	messageService         servicedDef.MessageService
	kafkaProducer          kafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) KafkaProducer(ctx context.Context) kafka.Producer {
	if d.kafkaProducer == nil {
		kafkaCfg := config.CommonAppConfig().Kafka
		syncProducer, err := sarama.NewSyncProducer(kafkaCfg.Brokers(), config.CommonAppConfig().Inbound.Config())
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return syncProducer.Close()
		})

		topics := []string{
			config.CommonAppConfig().Inbound.Topic(),
		}

		d.kafkaProducer = kafkaproducer.NewProducer(
			syncProducer, 
			topics, 
			logger.Logger(),
			platformKafkaMiddleware.TracingProducer(), // ← inject в headers
			platformKafkaMiddleware.LoggingProducer(logger.Logger()),
		)
	}
	return d.kafkaProducer
}

func (d *diContainer) MessageProducerService(ctx context.Context) servicedDef.ProducerService {
	if d.messageProducerService == nil {
		kafkaProducer := d.KafkaProducer(ctx)
		d.messageProducerService = inboundproducer.NewService(kafkaProducer)
	}
	return d.messageProducerService
}

func (d *diContainer) MessageService(ctx context.Context) servicedDef.MessageService {
	if d.messageService == nil {
		producerService := d.MessageProducerService(ctx)
		d.messageService = inbound.NewService(producerService)
	}
	return d.messageService
}

func (d *diContainer) MessageHandler(ctx context.Context) icv1.ServerInterface {
	if d.messageHandler == nil {
		messageService := d.MessageService(ctx)
		d.messageHandler = icv1.NewStrictHandler(api.New(messageService), nil)
	}
	return d.messageHandler
}
