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

func New(kind string, source sources.Source, counterType string) Processing {
	switch kind {
	case "simple":
		return NewSimple(source, counterType)
	case "advanced":
		return NewAdvanced(source, counterType)
	default:
		panic("unknown counter type")
	}
}