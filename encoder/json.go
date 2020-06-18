package encoder

import (
	"encoding/json"
	"github.com/slclub/link"
	"github.com/slclub/utils/bytesconv"
)

type Json struct {
}

func (js *Json) ContentType() string {
	return "application/json"
}

func (js *Json) Encode(data interface{}) string {
	data_bytes := js.EncodeBytes(data)
	return bytesconv.BytesToString(data_bytes)
}

func (js *Json) EncodeBytes(data interface{}) []byte {
	s, err := json.Marshal(data)
	if err != nil {
		link.ERROR("[JSON][ENCODER][convert string error]")
	}
	return (s)
}

func (js *Json) DecodeBytes(data []byte, obj interface{}) {
	err := json.Unmarshal(data, obj)
	if err != nil {
		link.ERROR("[JSON][DECODE][convert object]")
	}
}

func (js *Json) Decode(data string, obj interface{}) {
	js.DecodeBytes([]byte(data), obj)
}
