package sources

import (
	"net/http"
	"github.com/pkg/errors"
	"compress/gzip"
	"bufio"
	"log"
)

type Canned struct {
	content *gzip.Reader
	errors chan error
	records chan Record
}

func NewCanned() (Source, error) {
	resp, err := http.Get("http://tx.tamedia.ch.s3.amazonaws.com/challenge/data/stream.jsonl.gz")
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, errors.Errorf("HTTP error code %d", resp.StatusCode)
	}

	content, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	c := &Canned{content, make(chan error), make(chan Record)}
	go c.run()

	return c, nil
}

func (c *Canned) run() {
	scan := bufio.NewScanner(c.content)
	for scan.Scan() {
		line := scan.Bytes()
		rec, err := ParseRecord(line)
		if err != nil {
			c.errors <- err
		} else {
			c.records <- rec
		}
	}

	if err := scan.Err(); err != nil {
		log.Print(err)
	}
}

func (c *Canned) Errors() <-chan error {
	return c.errors
}

func (c *Canned) Records() <-chan Record {
	return c.records
}

func (c *Canned) Close() error {
	c.content.Close()

	return nil
}
