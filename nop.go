package log

import (
	"github.com/thecodinglab/log/level"
)

var _ Logger = (*nopLogger)(nil)

type nopLogger struct{}

func (l nopLogger) With(...interface{}) Logger     { return l }
func (l nopLogger) Level(level.Level) level.Logger { return nopLevelLogger{} }

func (l nopLogger) Debug(...interface{}) {}
func (l nopLogger) Info(...interface{})  {}
func (l nopLogger) Warn(...interface{})  {}
func (l nopLogger) Error(...interface{}) {}

func (l nopLogger) Debugf(string, ...interface{}) {}
func (l nopLogger) Infof(string, ...interface{})  {}
func (l nopLogger) Warnf(string, ...interface{})  {}
func (l nopLogger) Errorf(string, ...interface{}) {}

func (l nopLogger) Sync() {}

var _ level.Logger = (*nopLevelLogger)(nil)

type nopLevelLogger struct{}

func (l nopLevelLogger) With(...interface{}) level.Logger { return l }
func (l nopLevelLogger) Print(...interface{})             {}
func (l nopLevelLogger) Printf(string, ...interface{})    {}
