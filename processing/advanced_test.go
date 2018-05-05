package processing

import (
	"testing"
	"github.com/butzist/challenge/sources"
	"time"
	"fmt"
	"github.com/butzist/challenge/outputs"
)

func newTestee() (Processing, *mockSource, *mockClock) {
	s := &mockSource{make(chan error), make(chan sources.Record)}
	c := &mockClock{after: make(chan time.Time)}
	c.now = 6000000

	p := newAdvancedWithClock(s, c)
	p.(*Advanced).currentMinute = c.now / 60

	return p, s, c
}

func TestInOrder(t *testing.T) {
	p, s, c := newTestee()
	defer p.Close()

	s.records <- sources.Record{"a", 6000000}
	s.records <- sources.Record{"b", 6000059}
	s.records <- sources.Record{"c", 6000059}
	s.records <- sources.Record{"d", 6000060}

	time.Sleep(1)

	if c.now + uint64(c.afterTarget) != 6000065 {
		fmt.Errorf("wrong timeout")
	}

	c.now = 6000064
	c.after <- time.Now()
	assertNoOutputsReceived(p, t, "output too early")

	c.now = 6000065
	c.after <- time.Now()

	out1 := assertOutputReceived(p, t, "no output")
	if out1.Count != 3 {
		t.Errorf("count1 != 3")
	}

	c.now = 6000125
	c.after <- time.Now()

	out2 := assertOutputReceived(p, t, "no output")
	if out2.Count != 1 {
		t.Errorf("count2 != 1")
	}
}


func TestOutOfOrder(t *testing.T) {
	p, s, c := newTestee()
	defer p.Close()

	c.now = 6000064
	c.after <- time.Now()

	time.Sleep(1)

	s.records <- sources.Record{"a", 6000059}
	s.records <- sources.Record{"b", 6000058}

	time.Sleep(1)

	c.now = 6000065
	c.after <- time.Now()

	out := assertOutputReceived(p, t, "no output")
	if out.Count != 2 {
		t.Errorf("count != 2")
	}
}

func TestOutOfOrderError(t *testing.T) {
	p, s, c := newTestee()
	defer p.Close()

	c.now = 6000064
	c.after <- time.Now()

	time.Sleep(1)

	s.records <- sources.Record{"a", 6000000}
	s.records <- sources.Record{"b", 6000119}

	assertNoErrorsReceived(p, t, "valid timestamps not accepted")

	s.records <- sources.Record{"b", 6000120}
	assertErrorReceived(p, t, "no error received")

	s.records <- sources.Record{"b", 5999999}
	assertErrorReceived(p, t, "no error received")
}

func assertNoOutputsReceived(p Processing, t *testing.T, msg string) {
	time.Sleep(1)
	select {
	case <-p.Outputs():
		t.Error(msg)
	default:
		break
	}
}

func assertOutputReceived(p Processing, t *testing.T, msg string) *outputs.OutputStruct {
	time.Sleep(1)
	select {
	case out := <-p.Outputs():
		return out
	default:
		t.Error(msg)
	}
	return nil
}

func assertNoErrorsReceived(p Processing, t *testing.T, msg string) {
	time.Sleep(1)
	select {
	case <-p.Errors():
		t.Error(msg)
	default:
		break
	}
}

func assertErrorReceived(p Processing, t *testing.T, msg string) {
	time.Sleep(1)
	select {
	case <-p.Errors():
		return
	default:
		t.Error(msg)
	}
}

type mockSource struct {
	errors chan error
	records chan sources.Record
}

func (m *mockSource) Errors() <-chan error {
	return m.errors
}

func (m *mockSource) Records() <-chan sources.Record {
	return m.records
}

func (mockSource) Close() error {
	return nil
}

