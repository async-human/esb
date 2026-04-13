package kafka

import (
	"context"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg Message) error

// SendHandler — для middleware цепочки продюсера
type SendHandler func(ctx context.Context, msg Message) error

// ProducerMiddleware — аналог Middleware для консьюмера
type ProducerMiddleware func(next SendHandler) SendHandler

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, msg Message) error
}