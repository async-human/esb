package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	otlploggrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	log "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

var (
	globalLogger *logger
	once         sync.Once
	dynamicLevel zap.AtomicLevel // Позволяет изменять уровень логирования без перезапуска
)

// logger обёртка над zap.Logger с enrich поддержкой контекста
type logger struct {
	zapLogger *zap.Logger
}

type Config interface {
	Level() string
	AsJson() bool

	ServiceName() string
	Environment() string

	OTelEnabled() bool
	OTelCollectorEndpoint() string
}

// Init инициализирует глобальный логер с использованием паттерна Singleton
func Init(ctx context.Context, cfg Config) error {
	var err error
	once.Do(func() {
		// 1. Устанавливаем динамический уровень из конфигурации
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(cfg.Level()))

		// 2. Формируем набор ядер (Cores) для вывода логов
		cores, buildErr := buildCores(ctx, cfg)
		if buildErr != nil {
			err = buildErr
			return
		}

		// 3. Объединяем все источники (stdout + OTel) через NewTee
		combinedCore := zapcore.NewTee(cores...)

		// 4. Создаем экземпляр Zap логера с поддержкой информации о месте вызова
		zapLogger := zap.New(combinedCore,
			zap.AddCaller(),                   // Добавляем файл и строку вызова
			zap.AddCallerSkip(1),              // Пропускаем обертку пакета logger
			zap.AddStacktrace(zap.ErrorLevel), // Добавляем стектрейс для ERROR и выше
		)
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
	return err
}

// buildCores создает список ядер для записи в консоль и OpenTelemetry
func buildCores(ctx context.Context, cfg Config) ([]zapcore.Core, error) {
	var cores []zapcore.Core

	// Создаем ядро для стандартного вывода (stdout)
	cores = append(cores, createStdoutCore(cfg.AsJson(), dynamicLevel))

	// Если в конфиге включен OpenTelemetry, добавляем соответствующее ядро
	if cfg.OTelEnabled() {
		otelLogger, err := createOTLPLogger(ctx, cfg) // Метод инициализации OTel провайдера
		if err == nil {
			// Оборачиваем OTel логер в кастомный Core, реализующий интерфейс zapcore.Core
			cores = append(cores, NewSimpleOpenTelemetryCore(dynamicLevel, otelLogger))
		}
	}

	return cores, nil
}

func createOTLPLogger(ctx context.Context, cfg Config) (log.Logger, error) {
	// 1. Создаем OTLP gRPC экспортер для отправки данных в коллектор
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.OTelCollectorEndpoint()),
		otlploggrpc.WithInsecure(), // Отключаем TLS для локальной разработки
	)
	if err != nil {
		return nil, err
	}

	// 2. Описываем ресурсы сервиса (имя, версия, среда)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName()),
			attribute.String("environment", cfg.Environment()),
		),
	)
	if err != nil {
		return nil, err
	}

	// 3. Инициализируем LoggerProvider с пакетным процессором (BatchProcessor)
	// Пакетная отправка эффективнее для производительности приложения
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	// Возвращаем логер, привязанный к имени нашего сервиса
	return provider.Logger(cfg.ServiceName()), nil
}

// createStdoutCore конфигурирует стандартный вывод (JSON или текстовая консоль)
func createStdoutCore(asJSON bool, level zapcore.LevelEnabler) zapcore.Core {
	encoderConfig := buildEncoderConfig() // Базовые настройки полей

	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // Структурированный вид для продакшена
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Текстовый вид для разработки
	}

	// Возвращаем ядро, направленное в стандартный поток вывода
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// buildEncoderConfig определяет стандартные поля: время, уровень, сообщение
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO в верхнем регистре
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // Формат ISO8601
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // Краткий путь к файлу
	}
}

// SetLevel изменяет уровень логирования "на лету"
func SetLevel(level zapcore.Level) {
	dynamicLevel.SetLevel(level)
}

// GetLevel возвращает текущий уровень логирования
func GetLevel() zapcore.Level {
	return dynamicLevel.Level()
}

func InitForBenchmark(ctx context.Context, cfg Config) {
	Init(ctx, cfg)
}

// logger возвращает глобальный enrich-aware логгер
func Logger() *logger {
	return globalLogger
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с контекстом
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug enrich-aware debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

// Info enrich-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

// Warn enrich-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

// Error enrich-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

// Fatal enrich-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

// Instance methods для enrich loggers (logger)

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
