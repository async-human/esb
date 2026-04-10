package v1

import (
	"github.com/async-human/esb/inbound-connector/internal/service"
)

type api struct {
	messageService service.MessageService
}

func New(messageService service.MessageService) *api {
	return &api{
		messageService: messageService,
	}
}