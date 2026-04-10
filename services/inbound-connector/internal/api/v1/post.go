package v1

import (
	"context"

	"github.com/async-human/esb/inbound-connector/internal/converter"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
)

func (a *api) PostMessage(ctx context.Context, request icv1.PostMessageRequestObject) (icv1.PostMessageResponseObject, error) {
	ok, err := a.messageService.PostMessage(ctx, converter.MessageToModel(request))
	if err != nil {
		return nil, err
	}
	if !ok {
		return icv1.PostMessage400JSONResponse{}, nil
	}
	return icv1.PostMessage201Response{}, nil
}