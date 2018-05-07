package sources

import (
	"github.com/Shopify/sarama"
)

type Kafka struct {
	master sarama.Consumer
	consumer sarama.PartitionConsumer
	quit chan int
	errors chan error
	records chan *Record
}

func NewKafka(brokers []string, topic string, partition int32) (Source, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := master.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	k := &Kafka{master, consumer, make(chan int), make(chan error), make(chan *Record)}
	go k.run()

	return k, nil
}

func (k *Kafka) run() {
	for {
		select {
		case err := <-k.consumer.Errors():
			k.errors <- err
		case msg := <-k.consumer.Messages():
			rec, err := ParseRecord(msg.Value)
			if err != nil {
				k.errors <- err
			} else {
				k.records <- rec
			}
		case <-k.quit:
			return
		}
	}
}

func (k *Kafka) Errors() <-chan error {
	return k.errors
}

func (k *Kafka) Records() <-chan *Record {
	return k.records
}

func (k *Kafka) Close() error {
	close(k.quit)
	return k.master.Close()
}