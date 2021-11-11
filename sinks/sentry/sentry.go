package sentry

import (
	"errors"
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/thecodinglab/log"
	"github.com/thecodinglab/log/level"
)

var levelMapping = map[level.Level]sentry.Level{
	level.Debug: sentry.LevelDebug,
	level.Info:  sentry.LevelInfo,
	level.Warn:  sentry.LevelWarning,
	level.Error: sentry.LevelError,
}

type Option func(c *sink)

var _ log.Sink = (*sink)(nil)

type sink struct {
	min          level.Level
	hub          *sentry.Hub
	flushTimeout time.Duration
}

func New(opts ...Option) log.Sink {
	hub := sentry.CurrentHub()
	if hub.Client() == nil {
		client, err := sentry.NewClient(sentry.ClientOptions{})
		if err != nil {
			panic(err)
		}
		hub.BindClient(client)
	}

	s := sink{
		min:          level.Error,
		hub:          hub,
		flushTimeout: 5 * time.Second,
	}

	for _, option := range opts {
		option(&s)
	}

	return s
}

func (c sink) Write(entry *log.Entry) {
	if !c.min.Enabled(entry.Level) {
		return
	}

	event := sentry.NewEvent()

	for key, value := range entry.Fields {
		if err, ok := value.(error); ok {
			exceptions := extractExceptions(err)
			event.Exception = append(event.Exception, exceptions...)
			continue
		}

		event.Extra[key] = value
	}

	reverse(event.Exception)

	event.Level = levelMapping[entry.Level]
	event.Message = entry.Message
	event.Timestamp = entry.Time

	_ = c.hub.CaptureEvent(event)
}

func (c sink) Sync() {
	_ = c.hub.Flush(c.flushTimeout)
}

func WithLevel(min level.Level) Option {
	return func(c *sink) {
		c.min = min
	}
}

func WithHub(hub *sentry.Hub) Option {
	return func(c *sink) {
		c.hub = hub
	}
}

func extractExceptions(err error) []sentry.Exception {
	var exceptions []sentry.Exception

	for err != nil {
		stack := sentry.ExtractStacktrace(err)

		exceptions = append(exceptions, sentry.Exception{
			Type:       reflect.TypeOf(err).String(),
			Value:      err.Error(),
			Stacktrace: stack,
		})

		err = errors.Unwrap(err)
	}

	return exceptions
}

func reverse(a []sentry.Exception) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}
