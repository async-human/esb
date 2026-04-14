package v1

import (
	"context"
	"time"

	"github.com/async-human/esb/inbound-connector/internal/metrics"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func (a *api) GetInfo(ctx context.Context, request icv1.GetInfoRequestObject) (icv1.GetInfoResponseObject, error) {
	start := time.Now()

	info, err := a.messageService.GetInfo(ctx)

	statusCode := 200
	if err != nil {
		statusCode = 500
	}

	attrs := otelmetric.WithAttributes(
		attribute.String("http.request.method", "GET"),
		attribute.String("http.route", "/info"),
		attribute.Int("http.response.status_code", statusCode),
	)

	metrics.ServiceMetrics.HTTP.Duration.Record(ctx, time.Since(start).Seconds(), attrs)
	metrics.ServiceMetrics.HTTP.RequestsTotal.Add(ctx, 1, attrs)

	if err != nil {
		return nil, err
	}
	return icv1.GetInfo200JSONResponse{
		Name:    info.Name,
		Status:  info.Status,
		Version: info.Version,
	}, nil
}