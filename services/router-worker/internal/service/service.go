package service

import (
	"context"

	"github.com/async-human/esb/router-worker/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceMessage(ctx context.Context, message model.Message) error
}