package inboundconnector

import (
	"context"

	"github.com/async-human/esb/inbound-connector/internal/model"
	servicedDef "github.com/async-human/esb/inbound-connector/internal/service"
)

type service struct {
	messageProducerService servicedDef.MessageProducerService
}

func NewService(messageProducerService servicedDef.MessageProducerService) *service {
	return &service{
		messageProducerService: messageProducerService,
	}
}

func (s *service) GetInfo(ctx context.Context) (model.Info, error) {
	return model.Info{
		Name:    "inbound-connector",
		Status:  "ok",
		Version: "1.0.0",
	}, nil
}

func (s *service) PostMessage(ctx context.Context, request model.Message) (bool, error) {
	err := s.messageProducerService.ProduceMessageRecorded(ctx, request)
	if err != nil {
		return false, err
	}
	return true, nil
}
