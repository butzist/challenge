package main

import (
	"os"
	"os/signal"
	"log"
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
	"github.com/butzist/challenge/processing"
	"flag"
)

func main() {
	sourceType := flag.String("source", "kafka", "data source: kafka or canned")
	outputType := flag.String("output", "kafka", "data output: kafka or console")
	processingType := flag.String("processing", "advanced", "timestamp processing: simple or advanced")
	counterType := flag.String("counter", "exact", "cardinality counter: exact or probabilistic")
	flag.Parse()

	source, err := sources.New(*sourceType)
	if err != nil {
		panic(err)
	}
	defer source.Close()

	output, err := outputs.New(*outputType)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	process := processing.New(*processingType, source, *counterType)
	defer process.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-process.Errors():
			log.Println(err)
		case out := <-process.Outputs():
			err := output.Output(out)
			if err != nil {
				log.Println(err)
			}
		case <-signals:
			return
		}
	}
}
