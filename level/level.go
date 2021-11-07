package level

import "encoding/json"

var levelNames = map[Level]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
}

type Level int

const (
	Debug Level = 0
	Info  Level = 1
	Warn  Level = 2
	Error Level = 3
)

func (l Level) Enabled(level Level) bool {
	return level >= l
}

func (l Level) String() string {
	return levelNames[l]
}

func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

type Logger interface {
	With(kv ...interface{}) Logger

	Print(args ...interface{})
	Printf(format string, args ...interface{})
}
