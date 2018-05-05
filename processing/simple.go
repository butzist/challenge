package processing

import (
	"github.com/butzist/challenge/counters"
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
	"fmt"
)

type Simple struct {
	currentMinute uint64
	counter counters.Counter
}

func NewSimple() Processing {
	return &Simple{}
}

func (s *Simple) Process(record sources.Record) (*outputs.OutputStruct, error) {
	minute := uint64(record.Timestamp) / 60
	if s.currentMinute == minute {
		s.counter.Add(record.Uid)
	} else if s.currentMinute < minute {
		s.currentMinute = minute
		if s.counter != nil {
			out := &outputs.OutputStruct{s.counter.Count(), s.counter.Raw()}
			s.counter = counters.New()
			return out, nil
		} else {
			s.counter = counters.New()
		}
	} else {
		return nil, fmt.Errorf("out-of-order timestamp %d", uint64(record.Timestamp))
	}

	return nil, nil
}
