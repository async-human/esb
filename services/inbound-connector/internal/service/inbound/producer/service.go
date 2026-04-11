package producer

import (
	"github.com/async-human/esb/platform/kafka"
)

type service struct {
	producer kafka.Producer
}

func NewService(producer kafka.Producer) *service {
	return &service{
		producer: producer,
	}
}
