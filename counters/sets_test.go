package counters

import "testing"

func TestEmpty(t *testing.T) {
	s := NewSet()
	if s.Count() != 0 {
		t.Errorf("not empty")
	}
}

func TestUnique(t *testing.T) {
	s := NewSet()
	s.Add("a")
	s.Add("b")
	s.Add("c")

	if s.Count() != 3 {
		t.Errorf("not 3")
	}
}

func TestNotUnique(t *testing.T) {
	s := NewSet()
	s.Add("a")
	s.Add("a")
	s.Add("b")

	if s.Count() != 2 {
		t.Errorf("not 2")
	}
}

func TestAggregate(t *testing.T) {
	s1 := NewSet()
	s1.Add("a")
	s1.Add("b")

	s2 := NewSet()
	s2.Add("b")
	s2.Add("c")

	s1.Aggregate(s2)

	if s1.Count() != 3 {
		t.Errorf("aggregated not 3")
	}
}

func TestAggregateRaw(t *testing.T) {
	s1 := NewSet()
	s1.Add("a")
	s1.Add("b")

	s2 := NewSet()
	s2.Add("b")
	s2.Add("c")
	raw2 := s2.Raw()

	s1.AggregateRaw(raw2)

	if s1.Count() != 3 {
		t.Errorf("raw aggregated not 3")
	}
}