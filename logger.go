package log

import (
	"context"
	"time"

	"github.com/thecodinglab/log/fields"
	"github.com/thecodinglab/log/level"
)

const pkgName = "github.com/thecodinglab/log"

type Logger interface {
	With(kv ...interface{}) Logger

	Level(lvl level.Level) level.Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Sync()
}

type Entry struct {
	Level   level.Level `json:"level"`
	Time    time.Time   `json:"time"`
	Message string      `json:"message"`
	Caller  string      `json:"caller"`
	Fields  fields.KV   `json:"fields,omitempty"`
}

type Sink interface {
	Write(entry *Entry)
	Sync()
}

type loggerContextKey struct{}

func Get(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return logger
	}

	return nopLogger{}
}

func Attach(parent context.Context, logger Logger) context.Context {
	return context.WithValue(parent, loggerContextKey{}, logger)
}

func With(ctx context.Context, kv ...interface{}) Logger {
	return Get(ctx).With(kv...)
}

func Debug(ctx context.Context, args ...interface{}) {
	Get(ctx).Level(level.Debug).Print(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	Get(ctx).Level(level.Info).Print(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	Get(ctx).Level(level.Warn).Print(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	Get(ctx).Level(level.Error).Print(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	Get(ctx).Level(level.Debug).Printf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	Get(ctx).Level(level.Info).Printf(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	Get(ctx).Level(level.Warn).Printf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	Get(ctx).Level(level.Error).Printf(format, args...)
}

func Sync(ctx context.Context) {
	Get(ctx).Sync()
}
