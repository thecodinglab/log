package file

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/thecodinglab/log"
	"github.com/thecodinglab/log/fields"
	"github.com/thecodinglab/log/level"
)

var (
	resetColor  = color.New(color.Reset)
	levelColors = map[level.Level]*color.Color{
		level.Debug: color.New(color.FgMagenta),
		level.Info:  color.New(color.FgBlue),
		level.Warn:  color.New(color.FgYellow),
		level.Error: color.New(color.FgRed),
	}
)

type TextFormatter struct {
}

func (f TextFormatter) Format(buffer *bytes.Buffer, entry *log.Entry) error {
	buffer.WriteString(entry.Time.Format(time.RFC3339))
	buffer.WriteString("  ")

	buffer.WriteString(entry.Caller)
	buffer.WriteString("  ")

	lvl := fmt.Sprintf("%-5s", entry.Level.String())
	if col, ok := levelColors[entry.Level]; ok {
		_, _ = col.Fprintf(buffer, lvl)
		_, _ = resetColor.Fprint(buffer, "  ")
	} else {
		buffer.WriteString(lvl)
		buffer.WriteString("  ")
	}

	buffer.WriteString(entry.Message)

	if len(entry.Fields) != 0 {
		buffer.WriteString("  ")
		f.formatFields(buffer, entry.Fields, 0)
	}

	buffer.WriteRune('\n')
	return nil
}

func (f TextFormatter) formatFields(buffer *bytes.Buffer, kv fields.KV, depth int) {
	if depth > 3 {
		buffer.WriteString("{...}")
		return
	}

	if depth > 0 {
		buffer.WriteString("{")
	}

	first := true

	for k, v := range kv {
		if !first {
			buffer.WriteString(", ")
		}
		first = false

		buffer.WriteString(k)
		buffer.WriteRune('=')

		if nested, ok := v.(fields.KV); ok {
			f.formatFields(buffer, nested, depth+1)
		} else {
			buffer.WriteString(fmt.Sprint(v))
		}
	}

	if depth > 0 {
		buffer.WriteString("}")
	}
}
