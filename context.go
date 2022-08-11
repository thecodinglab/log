package log

import (
	"context"

	"github.com/thecodinglab/log/level"
)

type loggerContextKey struct{}

func Get(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return logger
	}

	return Default()
}

func Attach(parent context.Context, logger Logger) context.Context {
	return context.WithValue(parent, loggerContextKey{}, logger)
}

func With(ctx context.Context, kv ...interface{}) Logger {
	return Get(ctx).With(kv...)
}

func Level(ctx context.Context, lvl level.Level) LevelLogger {
	return Get(ctx).Level(lvl)
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
