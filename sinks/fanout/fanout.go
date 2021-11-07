package fanout

import (
	"sync"

	"github.com/thecodinglab/log"
)

var _ log.Sink = (*sink)(nil)

type sink struct {
	sinks []log.Sink
}

func New(destinations ...log.Sink) log.Sink {
	return sink{
		sinks: destinations,
	}
}

func (s sink) Write(entry *log.Entry) {
	s.forAll(func(sink log.Sink) {
		sink.Write(entry)
	})
}

func (s sink) Sync() {
	s.forAll(func(sink log.Sink) {
		sink.Sync()
	})
}

func (s sink) forAll(fn func(sink log.Sink)) {
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
