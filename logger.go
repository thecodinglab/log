package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/thecodinglab/log/fields"
	"github.com/thecodinglab/log/level"
)

const pkgName = "github.com/thecodinglab/log"

var defaultLogger = Logger{}

func Default() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
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

type Logger struct {
	sink Sink
	kv   fields.KV
}

func New(sink Sink) Logger {
	return Logger{
		sink: sink,
		kv:   nil,
	}
}

func (l Logger) IsValid() bool {
	return l.sink != nil
}

func (l Logger) With(kv ...interface{}) Logger {
	l.kv = fields.Merge(l.kv, fields.From(kv...))
	return l
}

func (l Logger) Level(lvl level.Level) LevelLogger {
	return LevelLogger{
		level:  lvl,
		sink:   l.sink,
		fields: l.kv,
	}
}

func (l Logger) Debug(args ...interface{}) {
	l.Level(level.Debug).Print(args...)
}

func (l Logger) Info(args ...interface{}) {
	l.Level(level.Info).Print(args...)
}

func (l Logger) Warn(args ...interface{}) {
	l.Level(level.Warn).Print(args...)
}

func (l Logger) Error(args ...interface{}) {
	l.Level(level.Error).Print(args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Level(level.Debug).Printf(format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Level(level.Info).Printf(format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.Level(level.Warn).Printf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Level(level.Error).Printf(format, args...)
}

func (l Logger) Sync() {
	if l.sink != nil {
		l.sink.Sync()
	}
}

type LevelLogger struct {
	level  level.Level
	sink   Sink
	fields fields.KV
}

func (l LevelLogger) With(kv ...interface{}) LevelLogger {
	l.fields = fields.Merge(l.fields, fields.From(kv...))
	return l
}

func (l LevelLogger) Print(args ...interface{}) {
	if l.sink == nil {
		return
	}

	l.write(fmt.Sprint(args...))
}

func (l LevelLogger) Printf(format string, args ...interface{}) {
	if l.sink == nil {
		return
	}

	l.write(fmt.Sprintf(format, args...))
}

func (l LevelLogger) write(message string) {
	caller := l.extractCaller(2)

	entry := &Entry{
		Level:   l.level,
		Time:    time.Now(),
		Message: message,
		Caller:  caller,
		Fields:  l.fields,
	}

	l.sink.Write(entry)
}

func (l LevelLogger) extractCaller(skip int) string {
	callers := make([]uintptr, 10)
	runtime.Callers(skip+2, callers)

	frames := runtime.CallersFrames(callers)
	for {
		frame, next := frames.Next()

		if !strings.HasPrefix(frame.Function, pkgName) {
			return fmt.Sprint(l.trimPath(frame.File), ":", frame.Line)
		}

		if !next {
			break
		}
	}

	return "<internal>"
}

func (l LevelLogger) trimPath(path string) string {
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	return path[idx+1:]
}
