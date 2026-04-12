package kafka

import "github.com/async-human/esb/outbound-connector/internal/model"

type MessageDecoder interface {
	Decode(data []byte) (model.Message, error)
}