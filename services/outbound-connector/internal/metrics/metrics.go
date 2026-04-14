package metrics

import (
	"errors"
	"fmt"

	"github.com/async-human/esb/outbound-connector/internal/config"
	shared "github.com/async-human/esb/platform/metrics"
	"go.opentelemetry.io/otel"
)

type Metrics struct {
	App      shared.App
	Consumer shared.KafkaConsumer
	Delivery shared.Delivery
}

var ServiceMetrics Metrics

func Init() error {
	meter := otel.Meter(config.CommonAppConfig().App.ServiceName())

	app,      e1 := shared.NewApp(meter)
	consumer, e2 := shared.NewKafkaConsumer(meter)
	delivery, e3 := shared.NewDelivery(meter)

	if err := errors.Join(e1, e2, e3); err != nil {
		return fmt.Errorf("outbound-connector metrics: %w", err)
	}

	ServiceMetrics = Metrics{app, consumer, delivery}
	return nil
}