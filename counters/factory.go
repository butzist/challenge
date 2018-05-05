package counters

type Counter interface {
	Add(key string)
	Aggregate(other Counter)
	AggregateRaw(interface{})
	Raw() interface{}
	Count() uint64
}

func New() Counter {
	return NewSet()
}