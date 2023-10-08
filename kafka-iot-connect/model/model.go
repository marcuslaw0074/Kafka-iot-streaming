package model

import (
	"encoding/json"
	"time"
)

func BMSRawDataEncoder(id string, status int8, value float64) BMSRawData {
	return BMSRawData{
		T:         "bms.raw",
		Value:     value,
		Timestamp: time.Now().UTC().UnixMilli(),
		Id:        id,
		Status:    status,
	}
}

type BMSRawData struct {
	T         string  `json:"_t"`
	Value     float64 `json:"value" format:"float64"`
	Timestamp int64   `json:"timestamp" example:"1" format:"int64"`
	Id        string  `json:"id"`
	Status    int8    `json:"status" example:"1" format:"int8"`
}

func (b BMSRawData) Encode() ([]byte, error) {
	return json.Marshal(b)
}

func (b BMSRawData) Length() int {
	bb, err := json.Marshal(b)
	if err != nil {
		return 0
	} else {
		return len(bb)
	}
}

type BMSDataType struct {
	T        string `json:"_t"`
	Id       string `json:"id"`
	ItemType string `json:"itemType"`
}

func BMSDataTypeEncoder(id, itemType string) BMSDataType {
	return BMSDataType{
		T:        "bms.type",
		Id:       id,
		ItemType: itemType,
	}
}


func (b BMSDataType) Encode() ([]byte, error) {
	return json.Marshal(b)
}

func (b BMSDataType) Length() int {
	bb, err := json.Marshal(b)
	if err != nil {
		return 0
	} else {
		return len(bb)
	}
}