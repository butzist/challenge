package main

import (
	"testing"
	"encoding/json"
	"github.com/butzist/challenge/sources"
	"math/rand"
	"strconv"
	"github.com/butzist/challenge/processing"
	"fmt"
	"time"
	"os"
	"runtime/pprof"
)

type benchmarkSource struct {
	records chan sources.Record
}

func (benchmarkSource) Errors() <-chan error {
	return make(chan error)
}

func (b *benchmarkSource) Records() <-chan sources.Record {
	return b.records
}

func (benchmarkSource) Close() error {
	return nil
}

func BenchmarkSetCounter(b *testing.B) {
	var lines [][]byte

	fmt.Printf("Creating %d records\n", b.N)
	dataLen := 0
	ts := uint64(time.Now().Unix())
	for i := 0; i < b.N; i++ {
		uid := strconv.Itoa(int(rand.Int63()))
		record := sources.Record{Timestamp:ts, Uid:uid}
		encoded, _ := json.Marshal(record)
		dataLen += len(encoded)
		lines = append(lines, encoded)
	}
	b.SetBytes(int64(dataLen/len(lines)))

	source := &benchmarkSource{make(chan sources.Record)}
	process := processing.New("simple", source, "exact")
	defer process.Close()
	quit := make(chan int)

	b.ResetTimer()

	fmt.Println("Started benchmark")

	f, err := os.Create("benchmark.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go func() {
		for _, msg := range lines {
			rec, _ := sources.ParseRecord(msg)
			source.records <- rec
		}
		close(quit)
	}()

	done := false
	for !done {
		select {
		case err := <-process.Errors():
			fmt.Println(err)
		case out := <-process.Outputs():
			fmt.Println(json.Marshal(out))
		case <- quit:
			done = true
			break
		}
	}

	last := sources.Record{Timestamp:ts+65, Uid:"last"}
	source.records <- last

	fmt.Println(<-process.Outputs())
}
