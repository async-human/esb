package kafka

import "github.com/async-human/esb/router-worker/internal/model"

type MessageDecoder interface {
	Decode(data []byte) (model.Message, error)
}