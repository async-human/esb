package app

import (
	"context"
	"fmt"

	"github.com/async-human/esb/inbound-connector/internal/config"
	"github.com/async-human/esb/inbound-connector/internal/metrics"
	"github.com/async-human/esb/platform/closer"
	"github.com/async-human/esb/platform/logger"
	metricsPlatform "github.com/async-human/esb/platform/metrics"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {

	logger.Info(ctx, "📊 Текущий уровень логирования: "+logger.GetLevel().String())
	logger.Debug(ctx, "🔍 [DEBUG] Inbound Connector: проверка уровня логирования DEBUG")
	logger.Info(ctx, "🚀 Inbound Connector запущен и готов к обработке сообщений")
	logger.Warn(ctx, "⚠️ [WARN] Inbound Connector: проверка уровня логирования WARN")
	logger.Error(ctx, "❌ [ERROR] Inbound Connector: проверка уровня логирования ERROR")

	metrics.AppStartsTotal.Add(ctx, 1)

	<-ctx.Done()

	metrics.AppEndTotal.Add(ctx, 1)

	logger.Info(ctx, "Shutdown signal received")

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initCloser,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {

	initLoggerConfig := struct {
		config.LoggerConfig
		config.OtelConfig
		config.AppConfig
	}{
		LoggerConfig: config.CommonAppConfig().Logger,
		OtelConfig:   config.CommonAppConfig().Otel,
		AppConfig:    config.CommonAppConfig().App,
	}

	return logger.Init(
		ctx,
		initLoggerConfig,
	)
}

func (a *App) initCloser(_ context.Context) error {
	
	closer.SetLogger(logger.Logger())
	closer.AddNamed("metrics", metricsPlatform.Shutdown)

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metricsPlatform.InitProvider(ctx, config.CommonAppConfig().MetricConfig)
	if err != nil {
		return fmt.Errorf("failed to init metrics: %w", err)
	}

	return nil
}
