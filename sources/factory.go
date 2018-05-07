package sources

import (
	"os"
	"strconv"
)

type Source interface {
	Errors() <-chan error
	Records() <-chan *Record
	Close() error
}

func New(kind string) (Source, error) {
	switch kind {
	case "canned":
		return NewCanned()
	case "kafka":
		partition, err := strconv.Atoi(getEnv("KAFKA_PARTITION", "0"))
		if err != nil {
			panic(err)
		}

		return NewKafka(
			[]string {getEnv("KAFKA_BROKER", "localhost:9092")},
			getEnv("KAFKA_TOPIC", "mytopic"),
			int32(partition),
		)
	default:
		panic("unknown counter type")
	}
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}