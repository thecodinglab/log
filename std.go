package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/thecodinglab/log/fields"
	"github.com/thecodinglab/log/level"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	sink Sink
	kv   fields.KV
}

func New(sink Sink) Logger {
	return stdLogger{
		sink: sink,
		kv:   nil,
	}
}

func (l stdLogger) With(kv ...interface{}) Logger {
	l.kv = fields.Merge(l.kv, fields.From(kv...))
	return l
}

func (l stdLogger) Level(lvl level.Level) level.Logger {
	return stdLevelLogger{
		level:  lvl,
		sink:   l.sink,
		fields: l.kv,
	}
}

func (l stdLogger) Debug(args ...interface{}) {
	l.Level(level.Debug).Print(args...)
}

func (l stdLogger) Info(args ...interface{}) {
	l.Level(level.Info).Print(args...)
}

func (l stdLogger) Warn(args ...interface{}) {
	l.Level(level.Warn).Print(args...)
}

func (l stdLogger) Error(args ...interface{}) {
	l.Level(level.Error).Print(args...)
}

func (l stdLogger) Debugf(format string, args ...interface{}) {
	l.Level(level.Debug).Printf(format, args...)
}

func (l stdLogger) Infof(format string, args ...interface{}) {
	l.Level(level.Info).Printf(format, args...)
}

func (l stdLogger) Warnf(format string, args ...interface{}) {
	l.Level(level.Warn).Printf(format, args...)
}

func (l stdLogger) Errorf(format string, args ...interface{}) {
	l.Level(level.Error).Printf(format, args...)
}

func (l stdLogger) Sync() {
	if l.sink != nil {
		l.sink.Sync()
	}
}

var _ level.Logger = (*stdLevelLogger)(nil)

type stdLevelLogger struct {
	level  level.Level
	sink   Sink
	fields fields.KV
}

func (l stdLevelLogger) With(kv ...interface{}) level.Logger {
	l.fields = fields.Merge(l.fields, fields.From(kv...))
	return l
}

func (l stdLevelLogger) Print(args ...interface{}) {
	if l.sink == nil {
		return
	}

	l.write(fmt.Sprint(args...))
}

func (l stdLevelLogger) Printf(format string, args ...interface{}) {
	if l.sink == nil {
		return
	}

	l.write(fmt.Sprintf(format, args...))
}

func (l stdLevelLogger) write(message string) {
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

func (l stdLevelLogger) extractCaller(skip int) string {
	callers := make([]uintptr, 10)
	runtime.Callers(skip+2, callers)

	frames := runtime.CallersFrames(callers)
	for {
		frame, next := frames.Next()

		if !strings.HasSuffix(frame.File, pkgName+"/std.go") {
			return fmt.Sprint(l.trimPath(frame.File), ":", frame.Line)
		}

		if !next {
			break
		}
	}

	return "<internal>"
}

func (l stdLevelLogger) trimPath(path string) string {
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
