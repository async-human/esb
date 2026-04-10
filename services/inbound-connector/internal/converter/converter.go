package converter

import (
	"github.com/async-human/esb/inbound-connector/internal/model"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func MessageToModel(message icv1.PostMessageRequestObject) model.Message {
	return model.Message{
		Id:      OpenAPIUUIDToUUID(message.Body.Id),
		Payload: message.Body.Payload,
	}
}

func InfoToAPI(info model.Info) icv1.GetInfo200JSONResponse {
	return icv1.GetInfo200JSONResponse{
		Name:    info.Name,
		Status:  info.Status,
		Version: info.Version,
	}
}

func OpenAPIUUIDToUUID(openapiUUID openapi_types.UUID) uuid.UUID {
	return uuid.UUID(openapiUUID)
}
