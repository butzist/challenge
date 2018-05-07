package sources

import "encoding/json"

type Record struct {
	Uid string `json:"uid"`
	Timestamp uint64 `json:"ts"`
}

func ParseRecord(raw []byte) (r *Record, err error) {
	r = &Record{}
	err = json.Unmarshal(raw, &r)
	return
}