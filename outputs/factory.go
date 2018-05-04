package outputs

type OutputStruct struct {
	Count uint64 `json:"count"`
	Raw interface{} `json:"raw"`
}

type Output interface {
	Output(out *OutputStruct) error
}

func New() (Output, error) {
	return NewConsole()
}