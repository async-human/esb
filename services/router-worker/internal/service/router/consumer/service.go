package consumer

import (
	"github.com/async-human/esb/platform/kafka"
	kafkaConverter "github.com/async-human/esb/router-worker/internal/converter/kafka"
	servicePkg "github.com/async-human/esb/router-worker/internal/service"
)

type service struct {
	consumer        kafka.Consumer
	producerService servicePkg.ProducerService
	consumerDecoder kafkaConverter.MessageDecoder
}

func NewService(consumer kafka.Consumer, producerService servicePkg.ProducerService, consumerDecoder kafkaConverter.MessageDecoder) *service {
	return &service{
		consumer:        consumer,
		producerService: producerService,
		consumerDecoder: consumerDecoder,
	}
}
