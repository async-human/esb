package metrics

import (
	"errors"
	"fmt"

	"github.com/async-human/esb/management-api/internal/config"
	shared "github.com/async-human/esb/platform/metrics"
	"go.opentelemetry.io/otel"
)

type Metrics struct {
	App      shared.App
	HTTP     shared.HTTP
	Producer shared.KafkaProducer
}

var ServiceMetrics Metrics

func Init() error {
	meter := otel.Meter(config.CommonAppConfig().App.ServiceName())

	app, e1 := shared.NewApp(meter)
	http, e2 := shared.NewHTTP(meter)
	producer, e3 := shared.NewKafkaProducer(meter)

	if err := errors.Join(e1, e2, e3); err != nil {
		return fmt.Errorf("management api metrics: %w", err)
	}

	ServiceMetrics = Metrics{app, http, producer}
	return nil
}