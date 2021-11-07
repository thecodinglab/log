package file

import (
	"bytes"
	"io"
	"sync"

	"github.com/thecodinglab/log"
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

var _ log.Sink = (*sink)(nil)

type sink struct {
	writer    WriteSyncer
	formatter Formatter
	pool      sync.Pool
}

func New(writer WriteSyncer, encoder Formatter) log.Sink {
	return &sink{
		writer:    writer,
		formatter: encoder,

		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (s *sink) Write(entry *log.Entry) {
	// TODO error handling
	_ = s.write(entry)
}

func (s *sink) Sync() {
	// TODO error handling
	_ = s.writer.Sync()
}

func (s *sink) write(entry *log.Entry) error {
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
