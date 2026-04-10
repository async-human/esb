package model

import "github.com/google/uuid"

type Info struct {
	Name    string
	Status  string
	Version string
}

type Message struct {
	Id      uuid.UUID
	Payload map[string]any
}
