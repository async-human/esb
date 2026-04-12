package service

import (
	"context"

	"github.com/async-human/esb/outbound-connector/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type SenderService interface {
	Send(ctx context.Context, message model.Message) error
}