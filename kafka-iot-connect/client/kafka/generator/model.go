package main

import (
	"encoding/json"
	"kafka-iot-connect/model"
)

type DataRow struct {
	T    string                 `json:"_t"`
	Data []model.TimeseriesData `json:"data"`
}

func (b DataRow) Encode() ([]byte, error) {
	return json.Marshal(b)
}

func (b DataRow) Length() int {
	bb, err := json.Marshal(b)
	if err != nil {
		return 0
	} else {
		return len(bb)
	}
}

type DataParser struct {
	Interval int       `json:"interval"`
	Data     []DataRow `json:"data"`
}
