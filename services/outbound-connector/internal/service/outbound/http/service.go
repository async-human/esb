package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
            Transport: otelhttp.NewTransport(http.DefaultTransport),
        },
    }
}