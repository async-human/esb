package service

import (
	"context"

	"github.com/async-human/esb/inbound-connector/internal/model"
)

type MessageService interface {
	GetInfo(ctx context.Context) (model.Info, error)
	PostMessage(ctx context.Context, request model.Message) (bool, error)
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type MessageProducerService interface {
	ProduceMessageRecorded(ctx context.Context, message model.Message) error
}