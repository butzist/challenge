package outputs

import "os"

type OutputStruct struct {
	Count uint64 `json:"count"`
	Timestamp uint64 `json:"ts"`
	Raw interface{} `json:"raw"`
}

type Output interface {
	Output(out *OutputStruct) error
	Close() error
}

func New(kind string) (Output, error) {
	switch kind {
	case "console":
		return NewConsole()
	case "kafka":
		return NewKafka(
			[]string {getEnv("KAFKA_BROKER", "localhost:9092")},
			getEnv("KAFKA_OUTPUT_TOPIC", "mycounts"),
		)
	default:
		panic("unknown output type")
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