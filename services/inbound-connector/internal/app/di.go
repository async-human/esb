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
	"github.com/async-human/esb/platform/kafka"
	kafkaproducer "github.com/async-human/esb/platform/kafka/producer"
	platformlogger "github.com/async-human/esb/platform/logger"
)

type diContainer struct {
	messageProducerService servicedDef.MessageProducerService
	messageHandler         icv1.ServerInterface
	messageService         servicedDef.MessageService
	kafkaProducer          kafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) KafkaProducer(ctx context.Context) (kafka.Producer, error) {
	if d.kafkaProducer == nil {
		kafkaCfg := config.CommonAppConfig().Kafka
		syncProducer, err := sarama.NewSyncProducer(kafkaCfg.Brokers(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka producer: %w", err)
		}

		d.kafkaProducer = kafkaproducer.NewProducer(syncProducer, "inbound-messages", platformlogger.Logger())
	}
	return d.kafkaProducer, nil
}

func (d *diContainer) MessageProducerService(ctx context.Context) (servicedDef.MessageProducerService, error) {
	if d.messageProducerService == nil {
		kafkaProducer, err := d.KafkaProducer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize kafka producer: %w", err)
		}
		d.messageProducerService = inboundproducer.NewService(kafkaProducer, "inbound-messages")
	}
	return d.messageProducerService, nil
}

func (d *diContainer) MessageService(ctx context.Context) (servicedDef.MessageService, error) {
	if d.messageService == nil {
		producerService, err := d.MessageProducerService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize message producer service: %w", err)
		}
		d.messageService = inbound.NewService(producerService)
	}
	return d.messageService, nil
}

func (d *diContainer) MessageHandler(ctx context.Context) (icv1.ServerInterface, error) {
	if d.messageHandler == nil {
		messageService, err := d.MessageService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize message service: %w", err)
		}
		d.messageHandler = icv1.NewStrictHandler(api.New(messageService), nil)
	}
	return d.messageHandler, nil
}
