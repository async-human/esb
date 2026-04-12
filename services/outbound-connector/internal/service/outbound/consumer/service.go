package consumer

import (
	"github.com/async-human/esb/platform/kafka"
	kafkaConverter "github.com/async-human/esb/outbound-connector/internal/converter/kafka"
	servicePkg "github.com/async-human/esb/outbound-connector/internal/service"
)

type service struct {
	consumer        kafka.Consumer
	endpoints       []string
	senderService   servicePkg.SenderService
	consumerDecoder kafkaConverter.MessageDecoder
}

func NewService(
	consumer kafka.Consumer,
	endpoints []string,
	senderService servicePkg.SenderService,
	consumerDecoder kafkaConverter.MessageDecoder,
) *service {
	return &service{
		consumer:        consumer,
		endpoints:       endpoints,
		senderService:   senderService,
		consumerDecoder: consumerDecoder,
	}
}
