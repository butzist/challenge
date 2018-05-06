package counters

import (
	"github.com/lytics/hll"
	"crypto/sha1"
	"encoding/binary"
	"encoding/base64"
)

type Hll struct {
	counter *hll.Hll
}

func NewHll() Counter {
	return &Hll{ hll.NewHll(14, 25)}
}

func hashU64(s string) uint64 {
	sha1Hash := sha1.Sum([]byte(s))
	return binary.LittleEndian.Uint64(sha1Hash[0:8])
}

func (h *Hll) Add(key string) {
	h.counter.Add(hashU64(key))
}

func (h *Hll) Aggregate(other Counter) {
	h.counter.Combine(other.(*Hll).counter)
}

func (h *Hll) AggregateRaw(other interface{}) {
	s := other.(string)
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	otherHll := hll.NewHll(14, 25)
	err = otherHll.UnmarshalPb(raw)
	if err != nil {
		panic(err)
	}

	h.counter.Combine(otherHll)
}

func (h *Hll) Raw() interface{} {
	raw, err := h.counter.MarshalPb()
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(raw)
}

func (h *Hll) Count() uint64 {
	return h.counter.Cardinality()
}
