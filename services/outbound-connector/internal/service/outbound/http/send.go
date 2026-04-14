package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/async-human/esb/outbound-connector/internal/config"
	"github.com/async-human/esb/outbound-connector/internal/metrics"
	"github.com/async-human/esb/outbound-connector/internal/model"
	"github.com/async-human/esb/platform/logger"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

func (s *service) Send(ctx context.Context, message model.Message) error {
	start := time.Now()

	endpoint := config.CommonAppConfig().InboundConsumerConfig.Endpoints()[0]

	payload, err := json.Marshal(message.Payload)
	if err != nil {
		logger.Error(ctx, "Failed to marshal payload", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error(ctx, "Failed to create request", zap.Error(err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)

	elapsed := time.Since(start).Seconds()

	// Базовые атрибуты, общие для всех метрик этого вызова
	baseAttrs := []attribute.KeyValue{
		attribute.String("server.address", endpoint),
		attribute.String("http.request.method", http.MethodPost),
	}

	if err != nil {
		metrics.ServiceMetrics.Delivery.AttemptsTotal.Add(ctx, 1,
			otelmetric.WithAttributes(
				append(baseAttrs, attribute.String("delivery.result", "connection_error"))...,
			),
		)
		metrics.ServiceMetrics.Delivery.Errors.Add(ctx, 1,
			otelmetric.WithAttributes(
				append(baseAttrs, attribute.String("delivery.error_type", classifyNetError(err)))...,
			),
		)
		// Exemplar — trace_id из span'а otelhttp.NewTransport
		metrics.ServiceMetrics.Delivery.Duration.Record(ctx, elapsed,
			otelmetric.WithAttributes(
				append(baseAttrs, attribute.String("delivery.result", "connection_error"))...,
			),
		)
		logger.Error(ctx, "Failed to send request", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	attrsWithStatus := append(baseAttrs,
		attribute.Int("http.response.status_code", res.StatusCode),
		attribute.String("delivery.result", classifyHTTPResult(res.StatusCode)),
	)

	// Exemplar прикрепится автоматически из активного span'а otelhttp.NewTransport
	metrics.ServiceMetrics.Delivery.Duration.Record(ctx, elapsed,
		otelmetric.WithAttributes(attrsWithStatus...),
	)
	metrics.ServiceMetrics.Delivery.AttemptsTotal.Add(ctx, 1,
		otelmetric.WithAttributes(attrsWithStatus...),
	)

	if res.StatusCode >= 400 {
		metrics.ServiceMetrics.Delivery.Errors.Add(ctx, 1,
			otelmetric.WithAttributes(attrsWithStatus...),
		)
		return fmt.Errorf("recipient returned %d", res.StatusCode)
	}

	logger.Info(ctx, "Request sent",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status_code", res.StatusCode),
	)

	return nil
}

// classifyHTTPResult — категория по status code, не raw число,
// чтобы не взрывать кардинальность метрик
func classifyHTTPResult(code int) string {
	switch {
	case code < 400:
		return "success"
	case code < 500:
		return "client_error"
	default:
		return "server_error"
	}
}

// classifyNetError — категория сетевой ошибки без raw строк в лейблах
func classifyNetError(err error) string {
	switch {
	case isTimeout(err):
		return "timeout"
	case isConnectionRefused(err):
		return "connection_refused"
	default:
		return "unknown"
	}
}

func isTimeout(err error) bool {
	if netErr, ok := err.(interface{ Timeout() bool }); ok {
		return netErr.Timeout()
	}
	return false
}

func isConnectionRefused(err error) bool {
	return err != nil && (containsString(err.Error(), "connection refused"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsString(s[1:], substr) ||
		s[:len(substr)] == substr)
}