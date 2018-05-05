package counters

type Sets struct {
	elements map[string]struct{}
}

func NewSet() Counter {
	return &Sets{make(map[string]struct{})}
}

func (m *Sets) Aggregate(other Counter) {
	otherSet := other.(*Sets)

	for key := range otherSet.elements {
		m.Add(key)
	}
}

func (m *Sets) AggregateRaw(other interface{}) {
	otherSlice := other.([]string)

	for _, key := range otherSlice {
		m.Add(key)
	}
}

func (m *Sets) Raw() interface{} {
	var raw []string
	for key := range m.elements {
		raw = append(raw, key)
	}

	return raw
}

func (m *Sets) Add(key string) {
	m.elements[key] = struct{}{}
}

func (m *Sets) Count() uint64 {
	return uint64(len(m.elements))
}