package outputs

type OutputStruct struct {
	Count uint64 `json:"count"`
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
		panic("NYI")
	default:
		panic("unknown output type")
	}}