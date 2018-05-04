package sources

import "encoding/json"

type Record struct {
	Uid string `json:"uid"`
	Timestamp float64 `json:"ts"`
}

func ParseRecord(raw []byte) (r Record, err error) {
	err = json.Unmarshal(raw, &r)
	return
}