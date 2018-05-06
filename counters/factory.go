package counters

type Counter interface {
	Add(key string)
	Aggregate(other Counter)
	AggregateRaw(interface{})
	Raw() interface{}
	Count() uint64
}

func New(kind string) Counter {
	switch kind {
	case "exact":
		return NewSet()
	case "probabilistic":
		panic("NYI")
	default:
		panic("unknown counter type")
	}
}