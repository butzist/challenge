package counters

type Counter interface {
	Add(key string)
	Count() uint64
}

func New() Counter {
	return &Mock{}
}