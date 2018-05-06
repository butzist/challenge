package outputs

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"strconv"
)

type Kafka struct {
	master sarama.SyncProducer
	topic string
}

func NewKafka(brokers []string, topic string) (Output, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner

	master, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	k := &Kafka{master, topic}
	return k, nil
}

func (k *Kafka) Output(out *OutputStruct) error {
	encoded, err := json.Marshal(out)
	if err != nil {
		return err
	}

	_,_,err = k.master.SendMessage(&sarama.ProducerMessage{
		Topic:k.topic,
		Key:sarama.StringEncoder(strconv.Itoa(int(out.Timestamp))),
		Value:sarama.StringEncoder(string(encoded)),
	})
	return err
}

func (k *Kafka) Close() error {
	return k.master.Close()
}