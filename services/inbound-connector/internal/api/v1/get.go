package v1

import (
	"context"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
)

func (a *api) GetInfo(ctx context.Context, request icv1.GetInfoRequestObject) (icv1.GetInfoResponseObject, error) {
	info, err := a.messageService.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return icv1.GetInfo200JSONResponse{
		Name:    info.Name,
		Status:  info.Status,
		Version: info.Version,
	}, nil
}
