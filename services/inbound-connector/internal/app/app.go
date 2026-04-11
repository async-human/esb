package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/async-human/esb/inbound-connector/internal/config"
	icv1 "github.com/async-human/esb/pkg/api/inbound-connector/v1"
	"github.com/async-human/esb/platform/closer"
	"github.com/async-human/esb/platform/logger"
	metricsPlatform "github.com/async-human/esb/platform/metrics"
	"github.com/async-human/esb/platform/tracing"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/swgui/v5emb"
	"go.uber.org/zap"
)

const (
	readHeaderTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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

	if err := a.runHttpServer(ctx); err != nil {
		logger.Error(ctx, "HTTP server crashed", zap.Error(err))
		return err
	}

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initHttpServer,
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

	initTracingCfg := struct {
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

func (a *App) initHttpServer(ctx context.Context) error {

	r := chi.NewRouter()

	// Добавляем отдачу самого OpenAPI JSON (сгенерирован автоматически из многофайловой ямлы!)
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		swagger, err := icv1.GetSwagger()
		if err != nil {
			http.Error(w, "Error loading swagger spec", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// swagger.MarshalJSON() выдаст склеенную, полноценную спеку
		jsonSpec, _ := swagger.MarshalJSON()
		w.Write(jsonSpec)
	})

	swaggerUI := v5emb.New("Inbound Connector API", "/swagger.json", "/docs")
	r.Mount("/docs", swaggerUI)

	handler := a.diContainer.MessageHandler(ctx)
	r.Mount("/", icv1.Handler(handler))

	a.httpServer = &http.Server{
		Addr:              config.CommonAppConfig().Rest.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			logger.Error(ctx, "HTTP server shutdown", zap.Error(err))
			return err
		}
		logger.Info(ctx, "✅ HTTP server stopped")
		return nil
	})

	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, "🚀 HTTP-сервер запущен %s\n", zap.String("addr", config.CommonAppConfig().Rest.Address()))
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Ошибка запуска сервера", zap.Error(err))
		return err
	}
	return nil
}
