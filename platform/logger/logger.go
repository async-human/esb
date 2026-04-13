package logger

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/attribute"
	otlploggrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	userIDKey Key = "user_id"
)

var (
	globalLogger   *logger
	once           sync.Once
	dynamicLevel   zap.AtomicLevel
	globalProvider *sdklog.LoggerProvider
	providerMu     sync.RWMutex
)

type logger struct {
	log *zap.Logger
}

type Config interface {
	Level() string
	AsJson() bool

	ServiceName() string
	Environment() string

	OTelEnabled() bool
	OTelCollectorEndpoint() string
}

func Init(ctx context.Context, cfg Config) error {
	var err error
	once.Do(func() {
		parsedLevel, parseErr := zapcore.ParseLevel(cfg.Level())
		if parseErr != nil {
			parsedLevel = zapcore.InfoLevel
		}
		dynamicLevel = zap.NewAtomicLevelAt(parsedLevel)

		// 1. Создаем базовое ядро (stdout)
		stdoutCore := createStdoutCore(cfg.AsJson(), dynamicLevel)

		var finalCore zapcore.Core

		// 2. Если OTel включен, объединяем ядра через Tee
		if cfg.OTelEnabled() {
			provider, initErr := createOTLPProvider(ctx, cfg)
			if initErr != nil {
				err = initErr
				return
			}
			global.SetLoggerProvider(provider)

			// Официальный мост создает zapcore.Core
			otelCore := otelzap.NewCore(
				cfg.ServiceName(), 
				otelzap.WithLoggerProvider(provider),
			)
			
			// Направляем логи сразу в два потока: в консоль и в OTel
			finalCore = zapcore.NewTee(stdoutCore, otelCore)
		} else {
			finalCore = stdoutCore
		}

		// 3. Создаем финальный логгер
		finalZapLogger := zap.New(finalCore,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.ErrorLevel),
		)

		globalLogger = &logger{
			log: finalZapLogger,
		}
	})
	return err
}

func createOTLPProvider(ctx context.Context, cfg Config) (*sdklog.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.OTelCollectorEndpoint()),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName()),
			attribute.String("environment", cfg.Environment()),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	providerMu.Lock()
	globalProvider = provider
	providerMu.Unlock()

	return provider, nil
}

func createStdoutCore(asJSON bool, level zapcore.LevelEnabler) zapcore.Core {
	encoderConfig := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func SetLevel(level zapcore.Level) {
	dynamicLevel.SetLevel(level)
}

func GetLevel() zapcore.Level {
	return dynamicLevel.Level()
}

func Shutdown(ctx context.Context) error {
	providerMu.RLock()
	provider := globalProvider
	providerMu.RUnlock()

	if provider == nil {
		return nil
	}

	return provider.Shutdown(ctx)
}

func Logger() *logger {
	return globalLogger
}

func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{log: zap.NewNop()}
	}

	return &logger{
		log: globalLogger.log.With(fields...),
	}
}

func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{log: zap.NewNop()}
	}
	
	fields := fieldsFromContext(ctx)
	return With(fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Fatal(ctx, msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.log.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.log.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.log.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.log.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.log.Fatal(msg, allFields...)
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	fields := make([]zap.Field, 0, 3)

	// 1. Достаем данные трейсинга из контекста
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		fields = append(fields, zap.String("trace_id", spanContext.TraceID().String()))
	}
	if spanContext.HasSpanID() {
		fields = append(fields, zap.String("span_id", spanContext.SpanID().String()))
	}

	// 2. Достаем бизнес-данные
	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}