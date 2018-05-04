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

	output, err := outputs.New()
	if err != nil {
		panic(err)
	}

	processing := processing.New()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-source.Errors():
			log.Println(err)
		case rec := <-source.Records():
			if out, err := processing.Process(rec); out != nil {
				err2 := output.Output(out)
				if err2 != nil {
					log.Println(err2)
				}
			} else if err != nil {
				log.Println(err)
			}
		case <-signals:
			return
		}
	}
}
