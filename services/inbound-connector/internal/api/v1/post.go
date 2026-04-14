package v1

import (
	"context"
	"time"

	"github.com/async-human/esb/inbound-connector/internal/converter"
	"github.com/async-human/esb/inbound-connector/internal/metrics"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func (a *api) PostMessage(ctx context.Context, request icv1.PostMessageRequestObject) (icv1.PostMessageResponseObject, error) {
	start := time.Now()

	// InFlight: сколько PostMessage выполняется прямо сейчас
	metrics.ServiceMetrics.App.InFlight.Add(ctx, 1)
	defer metrics.ServiceMetrics.App.InFlight.Add(ctx, -1)

	ok, err := a.messageService.PostMessage(ctx, converter.MessageToModel(request))

	// Определяем статус до записи метрик — чтобы attrs были точными
	statusCode := 201
	switch {
	case err != nil:
		statusCode = 500
	case !ok:
		statusCode = 400
	}

	attrs := otelmetric.WithAttributes(
		attribute.String("http.request.method", "POST"),
		attribute.String("http.route", "/message"),
		attribute.Int("http.response.status_code", statusCode),
	)

	// ctx содержит активный span от otelhttp — SDK автоматически
	// прикрепит trace_id как Exemplar к этой histogram записи
	metrics.ServiceMetrics.HTTP.Duration.Record(ctx, time.Since(start).Seconds(), attrs)
	metrics.ServiceMetrics.HTTP.RequestsTotal.Add(ctx, 1, attrs)

	if err != nil {
		return nil, err
	}
	if !ok {
		return icv1.PostMessage400JSONResponse{}, nil
	}
	return icv1.PostMessage201Response{}, nil
}