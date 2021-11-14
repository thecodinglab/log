package file

import (
	"bytes"
	"io"
	"sync"

	"github.com/thecodinglab/log"
	"github.com/thecodinglab/log/level"
)

type WriteSyncer interface {
	io.Writer

	Sync() error
}

type Formatter interface {
	Format(buffer *bytes.Buffer, entry *log.Entry) error
}

type FormatterFunc func(buffer *bytes.Buffer, entry *log.Entry) error

func (fn FormatterFunc) Format(buffer *bytes.Buffer, entry *log.Entry) error {
	return fn(buffer, entry)
}

var _ log.Sink = (*Sink)(nil)

type Sink struct {
	writer    WriteSyncer
	formatter Formatter
	min       level.Level
	pool      sync.Pool
}

func New(writer WriteSyncer, encoder Formatter, min level.Level) log.Sink {
	return &Sink{
		writer:    writer,
		formatter: encoder,
		min:       min,

		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (s *Sink) SetMinLevel(min level.Level) {
	s.min = min
}

func (s *Sink) Write(entry *log.Entry) {
	if !s.min.Enabled(entry.Level) {
		return
	}

	// TODO error handling
	_ = s.write(entry)
}

func (s *Sink) Sync() {
	// TODO error handling
	_ = s.writer.Sync()
}

func (s *Sink) write(entry *log.Entry) error {
	buffer := s.pool.Get().(*bytes.Buffer)
	defer s.pool.Put(buffer)

	buffer.Reset()

	if err := s.formatter.Format(buffer, entry); err != nil {
		return err
	}

	if _, err := buffer.WriteTo(s.writer); err != nil {
		return err
	}

	return nil
}
