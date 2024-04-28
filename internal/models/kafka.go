package models

import "encoding/json"

type SearchMsg struct {
	TgId int
}

func (msg SearchMsg) Bytes() ([]byte, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return bytes, err
}
