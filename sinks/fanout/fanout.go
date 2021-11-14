package fanout

import (
	"sync"

	"github.com/thecodinglab/log"
)

var _ log.Sink = (*Sink)(nil)

type Sink struct {
	sinks []log.Sink
}

func New(destinations ...log.Sink) *Sink {
	return &Sink{
		sinks: destinations,
	}
}

func (s *Sink) Append(sinks ...log.Sink) {
	s.sinks = append(s.sinks, sinks...)
}

func (s *Sink) Write(entry *log.Entry) {
	s.forAll(func(sink log.Sink) {
		sink.Write(entry)
	})
}

func (s *Sink) Sync() {
	s.forAll(func(sink log.Sink) {
		sink.Sync()
	})
}

func (s *Sink) forAll(fn func(sink log.Sink)) {
	wg := sync.WaitGroup{}
	wg.Add(len(s.sinks))

	for _, sink := range s.sinks {
		go func(sink log.Sink) {
			defer wg.Done()
			fn(sink)
		}(sink)
	}

	wg.Wait()
}
