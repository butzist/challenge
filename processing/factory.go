package processing

import (
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
)

type Processing interface {
	Process(record sources.Record) (*outputs.OutputStruct, error)
}

func New() Processing {
	return NewSimple()
}