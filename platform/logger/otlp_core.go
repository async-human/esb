package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

// SimpleOpenTelemetryCore реализует интерфейс zapcore.Core
type SimpleOpenTelemetryCore struct {
	level  zapcore.LevelEnabler // Уровень логирования
	logger log.Logger           // Логер из пакета OpenTelemetry
}

// NewSimpleOpenTelemetryCore создает новый экземпляр Core для OTel
func NewSimpleOpenTelemetryCore(level zapcore.LevelEnabler, logger log.Logger) *SimpleOpenTelemetryCore {
	return &SimpleOpenTelemetryCore{
		level:  level,
		logger: logger,
	}
}

// Enabled проверяет, должен ли лог данного уровня быть записан
func (c *SimpleOpenTelemetryCore) Enabled(lvl zapcore.Level) bool {
	return c.level.Enabled(lvl)
}

// With создает копию Core с дополнительными полями
func (c *SimpleOpenTelemetryCore) With(fields []zapcore.Field) zapcore.Core {
	// В данной реализации поля обрабатываются непосредственно в методе Write
	return &SimpleOpenTelemetryCore{
		level:  c.level,
		logger: c.logger,
	}
}

// Check добавляет этот Core в CheckedEntry, если уровень логирования проходит проверку
func (c *SimpleOpenTelemetryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write преобразует запись Zap в формат OpenTelemetry Record и отправляет её
func (c *SimpleOpenTelemetryCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {

	record := log.Record{}
	record.SetTimestamp(ent.Time)                // Установка времени из Zap
	record.SetBody(log.StringValue(ent.Message)) // Тест сообщения

	// Мапим уровень логирования Zap на Severity OpenTelemetry
	severity := c.mapSeverity(ent.Level)
	record.SetSeverity(severity)

	// Добавляем caller-информацию (файл:строка)
	if ent.Caller.Defined {
		record.AddAttributes(
			log.String("caller", ent.Caller.String()),
			log.String("function", ent.Caller.Function),
		)
	}

	// Добавляем стектрейс для WARN и выше
	if ent.Stack != "" {
		record.AddAttributes(log.String("stacktrace", ent.Stack))
	}

	// Преобразуем атрибуты (поля) лога Zap в KeyValue OTel
	attrs := c.encodeFieldsToAttrs(fields, severity)
	if len(attrs) > 0 {
		record.AddAttributes(attrs...)
	}

	// Отправляем лог в коллектор с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	c.logger.Emit(ctx, record)

	return nil
}

// Sync является обязательным методом интерфейса, но в данной реализации ничего не делает
func (c *SimpleOpenTelemetryCore) Sync() error {
	return nil
}

// mapSeverity преобразует уровни Zap в типы Severity OpenTelemetry
func (c *SimpleOpenTelemetryCore) mapSeverity(lvl zapcore.Level) log.Severity {
	switch lvl {
	case zapcore.DebugLevel:
		return log.SeverityDebug
	case zapcore.InfoLevel:
		return log.SeverityInfo
	case zapcore.WarnLevel:
		return log.SeverityWarn
	case zapcore.ErrorLevel:
		return log.SeverityError
	case zapcore.FatalLevel:
		return log.SeverityFatal
	default:
		return log.SeverityUndefined
	}
}

func (c *SimpleOpenTelemetryCore) encodeFieldsToAttrs(fields []zapcore.Field, severity log.Severity) []log.KeyValue {

	if len(fields) == 0 && severity == log.SeverityUndefined {
		return nil
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	attrs := make([]log.KeyValue, 0, len(enc.Fields)+1)
	attrs = append(attrs, log.String("severity", severity.String()))
	for k, v := range enc.Fields {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, log.String(k, val))
		case bool:
			attrs = append(attrs, log.Bool(k, val))
		case int64:
			attrs = append(attrs, log.Int64(k, val))
		case float64:
			attrs = append(attrs, log.Float64(k, val))
		}
	}

	return attrs

}
