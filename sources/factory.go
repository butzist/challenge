package sources

type Source interface {
	Errors() <-chan error
	Records() <-chan Record
	Close() error
}

func New() (Source, error) {
	return NewCanned()
}