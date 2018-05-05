package outputs

import (
	"encoding/json"
	"os"
)

type Console struct {}

func NewConsole() (Output, error) {
	return &Console{}, nil
}

func (*Console) Output(out *OutputStruct) error {
	encoded, err := json.Marshal(out)
	if err != nil {
		return err
	}

	os.Stdout.Write(encoded)
	os.Stdout.Write([]byte{'\n'})

	return nil
}

func (*Console) Close() error {
	return nil
}