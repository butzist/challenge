package processing

import (
	"github.com/butzist/challenge/counters"
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
	"fmt"
	"time"
)

// Advanced timestamp processing allows records to arrive up to 5s late
// and still be processed within the correct counting bucket. Depends
// on consumer time being synced with producer time.
type Advanced struct {
	currentMinute uint64
	counter counters.Counter
	nextCounter counters.Counter
	errors chan error
	outputs chan *outputs.OutputStruct
	quit chan int
	clock clock
}

func NewAdvanced(source sources.Source) Processing {
	s := newAdvancedWithClock(source, &realClock{})
	return s
}

func newAdvancedWithClock(source sources.Source, clock clock) Processing {
	s := &Advanced{uint64(time.Now().Unix() / 60), counters.New(), counters.New(), make(chan error), make(chan *outputs.OutputStruct), make(chan int), clock}
	go s.run(source)
	return s
}

func (s *Advanced) run(source sources.Source) {
	for {
		currentTime := s.clock.Now()
		currentMinute := (currentTime - 5) / 60

		// send data for last minute if 5 seconds after next minute passed
		if currentMinute > s.currentMinute {
			s.currentMinute = currentTime / 60
			out := &outputs.OutputStruct{Count: s.counter.Count(), Raw: s.counter.Raw()}
			s.nextCounter,s.counter = counters.New(),s.nextCounter
			s.outputs <- out
		}

		select {
		case err := <-source.Errors():
			s.errors <- err
		case record := <-source.Records():
			minute := uint64(record.Timestamp) / 60
			if s.currentMinute == minute {
				s.counter.Add(record.Uid)
			} else if s.currentMinute + 1 == minute {
				s.nextCounter.Add(record.Uid)
			} else {
				s.errors <- fmt.Errorf("out-of-order timestamp %d", uint64(record.Timestamp))
			}
		case <-s.clock.After(time.Duration(s.currentMinute * 60 + 65 - currentTime) * time.Second):
			break
		case <-s.quit:
			return
		}
	}
}

func (s *Advanced) Errors() <-chan error {
	return s.errors
}

func (s *Advanced) Outputs() <-chan *outputs.OutputStruct {
	return s.outputs
}

func (s *Advanced) Close() error {
	close(s.quit)
	return nil
}