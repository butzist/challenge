package main

import (
	"os"
	"os/signal"
	"log"
	"github.com/butzist/challenge/sources"
	"github.com/butzist/challenge/outputs"
	"github.com/butzist/challenge/processing"
)

func main() {
	source, err := sources.New()
	if err != nil {
		panic(err)
	}
	defer source.Close()

	output, err := outputs.New()
	if err != nil {
		panic(err)
	}
	defer output.Close()

	process := processing.New(source)
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
