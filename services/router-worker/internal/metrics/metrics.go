package metrics

import (
	"errors"
	"fmt"

	shared "github.com/async-human/esb/platform/metrics"
	"github.com/async-human/esb/router-worker/internal/config"
	"go.opentelemetry.io/otel"
)

type Metrics struct {
	App      shared.App
	Producer shared.KafkaProducer
	Routing  shared.Routing
}

var ServiceMetrics Metrics

func Init() error {
	meter := otel.Meter(config.CommonAppConfig().App.ServiceName())

	app,      e1 := shared.NewApp(meter)
	producer, e3 := shared.NewKafkaProducer(meter)
	routing,  e4 := shared.NewRouting(meter)

	if err := errors.Join(e1, e3, e4); err != nil {
		return fmt.Errorf("router-worker metrics: %w", err)
	}

	ServiceMetrics = Metrics{app, producer, routing}
	return nil
}