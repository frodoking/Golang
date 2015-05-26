package gof

import (
	"encoding/json"
)

//操作Json结果
type Message struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

//序列化
func (r Message) Marshal() []byte {
	json, _ := json.Marshal(r)
	return json
}
