package processing

import (
	"time"
)

// Abstraction of time and timeout mechanism for unit testing
type clock interface {
	Now() uint64
	After(duration time.Duration) <-chan time.Time
}

type realClock struct {}

func (realClock) Now() uint64 {
	return uint64(time.Now().Unix())
}

func (realClock) After(duration time.Duration) <-chan time.Time {
	return time.After(duration)
}

type mockClock struct {
	now uint64
	after chan time.Time
	afterTarget time.Duration
}

func (m *mockClock) Now() uint64 {
	return m.now
}

func (m *mockClock) After(duration time.Duration) <-chan time.Time {
	m.afterTarget = duration
	return m.after
}
