package capture

import (
	"context"
	"errors"

	"github.com/thecodinglab/log"
	"github.com/thecodinglab/log/level"
)

const ErrorField = "error"

type Option func(config *Config)

type Config struct {
	Level level.Level
}

func Error(ctx context.Context, err error, opts ...Option) {
	if err == nil {
		return
	}

	config := Config{
		Level: level.Error,
	}

	for _, opt := range opts {
		opt(&config)
	}

	ErrorWithConfig(ctx, config, err)
}

func ErrorWithConfig(ctx context.Context, config Config, err error) {
	if err == nil {
		return
	}

	// check if we have multiple errors and if so separate them into their individual errors
	var group interface{ Errors() []error }
	if errors.As(err, &group) {
		for _, child := range group.Errors() {
			ErrorWithConfig(ctx, config, child)
		}
		return
	}

	log.Get(ctx).
		Level(config.Level).
		With(ErrorField, err).
		Print(err.Error())
}

func WithLevel(level level.Level) Option {
	return func(config *Config) {
		config.Level = level
	}
}
