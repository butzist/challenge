package processing

import (
	"github.com/butzist/challenge/counters"
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
	"fmt"
)

// Simple timestamp processing requires strictly ordered timestamps but on the other hand does not
// require system time to be synced with the producer i.e. can be used for importing historical data.
type Simple struct {
	currentMinute uint64
	counter counters.Counter
	errors chan error
	outputs chan *outputs.OutputStruct
	quit chan int
	counterType string
}

func NewSimple(source sources.Source, counterType string) Processing {
	s := &Simple{0, nil, make(chan error), make(chan *outputs.OutputStruct), make(chan int), counterType}
	go s.run(source)
	return s
}

func (s *Simple) run(source sources.Source) {
	for {
		select {
		case err := <-source.Errors():
			s.errors <- err
		case record := <-source.Records():
			minute := uint64(record.Timestamp) / 60
			if s.currentMinute == minute {
				s.counter.Add(record.Uid)
			} else if s.currentMinute < minute {
				s.currentMinute = minute
				if s.counter != nil {
					out := &outputs.OutputStruct{Count: s.counter.Count(), Raw: s.counter.Raw()}
					s.counter = counters.New(s.counterType)
					s.outputs <- out
				} else {
					s.counter = counters.New(s.counterType)
				}
			} else {
				s.errors <- fmt.Errorf("out-of-order timestamp %d", uint64(record.Timestamp))
			}
		case <-s.quit:
			return

		}
	}
}

func (s *Simple) Errors() <-chan error {
	return s.errors
}

func (s *Simple) Outputs() <-chan *outputs.OutputStruct {
	return s.outputs
}

func (s *Simple) Close() error {
	close(s.quit)
	return nil
}