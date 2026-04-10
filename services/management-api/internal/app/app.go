package app

import (
	"context"
	"fmt"

	"github.com/async-human/esb/management-api/internal/config"
	"github.com/async-human/esb/management-api/internal/metrics"
	"github.com/async-human/esb/platform/closer"
	"github.com/async-human/esb/platform/logger"
	metricsPlatform "github.com/async-human/esb/platform/metrics"
	"github.com/async-human/esb/platform/tracing"
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
	logger.Debug(ctx, "🔍 [DEBUG] Management API: проверка уровня логирования DEBUG")
	logger.Info(ctx, "🚀 Management API запущен и готов к работе")
	logger.Warn(ctx, "⚠️ [WARN] Management API: проверка уровня логирования WARN")
	logger.Error(ctx, "❌ [ERROR] Management API: проверка уровня логирования ERROR")

	metrics.AppStartsTotal.Add(ctx, 1)
	_, span := tracing.StartSpan(ctx, "start.management-api")

	<-ctx.Done()

	span.End()
	metrics.AppEndTotal.Add(ctx, 1)

	logger.Info(ctx, "Shutdown signal received")

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
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

	loggerItem := logger.Init(
		ctx,
		initLoggerConfig,
	)

	closer.SetLogger(logger.Logger())
	closer.AddNamed("logger", logger.Shutdown)

	return loggerItem
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metricsPlatform.InitProvider(ctx, config.CommonAppConfig().MetricConfig)
	if err != nil {
		return fmt.Errorf("failed to init metrics: %w", err)
	}

	closer.AddNamed("metrics", metricsPlatform.Shutdown)

	return nil
}

func (a *App) initTracing(ctx context.Context) error {

	initTracingCfg := struct{
		config.OtelConfig
		config.AppConfig
	}{
		OtelConfig: config.CommonAppConfig().Otel,
		AppConfig:  config.CommonAppConfig().App,
	}

	err := tracing.InitTracer(ctx, initTracingCfg)
	if err != nil {
		return err
	}

	closer.AddNamed("tracing", tracing.ShutdownTracer)

	return nil
}

