package main

import (
	"encoding/json"
)

func marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func unmarshal(data []byte, to interface{}) (interface{}, error) {
	err := json.Unmarshal(data, to)
	if err != nil {
		return nil, err
	}
	return to, nil
}