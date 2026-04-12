package http

import (
	"net/http"
	"time"
)

const(
	timeout = time.Second * 10
)

type service struct {
	client *http.Client
}

func NewService() *service {
	return &service{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}