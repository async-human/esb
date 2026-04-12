package model

import "github.com/google/uuid"

type Message struct {
	Id      uuid.UUID
	Payload map[string]any
}