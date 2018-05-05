package processing

import (
	"github.com/butzist/challenge/outputs"
	"github.com/butzist/challenge/sources"
)

type Processing interface {
	Errors() <-chan error
	Outputs() <-chan *outputs.OutputStruct
	Close() error
}

func New(source sources.Source) Processing {
	return NewSimple(source)
}