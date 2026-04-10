package producer

import (
	"github.com/async-human/esb/platform/kafka"
)

type service struct {
	producer kafka.Producer
	topic    string
}

func NewService(producer kafka.Producer, topic string) *service {
	return &service{
		producer: producer,
		topic:    topic,
	}
}
