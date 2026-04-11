package decode

import (
	"encoding/json"

	"github.com/async-human/esb/router-worker/internal/model"
)

type decoder struct{}

func NewDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.Message, error) {
	var message model.Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}