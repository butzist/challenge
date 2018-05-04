package counters

type Mock struct {
	count uint64
}

func (m *Mock) Add(key string) {
	m.count++
}

func (m *Mock) Count() uint64 {
	return m.count
}