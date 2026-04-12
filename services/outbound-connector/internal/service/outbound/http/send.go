package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/async-human/esb/outbound-connector/internal/config"
	"github.com/async-human/esb/outbound-connector/internal/model"
	"github.com/async-human/esb/platform/logger"
	"go.uber.org/zap"
)

func (s *service) Send(ctx context.Context, message model.Message) error {

	payload, err := json.Marshal(message.Payload)
	if err != nil {
		logger.Error(ctx, "Failed to marshal payload", zap.Error(err))
		return err
	}

	endpoint := config.CommonAppConfig().InboundConsumerConfig.Endpoints()[0]

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error(ctx, "Failed to create request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		logger.Error(ctx, "Failed to send request", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	logger.Info(ctx, "Request sent",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status_code", res.StatusCode),
	)

	return nil
}
