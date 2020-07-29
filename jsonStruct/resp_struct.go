package jsonStruct

import (
	"encoding/json"
	"log"
)

// 基本上都照着这个来做接口的返回值
type CommonResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func CommonResp(c CommonResponse) ([]byte, error){
	v, err := json.Marshal(c)
	if err!=nil {
		log.Println(err)
		return []byte("a"), err
	}
	return v, nil
}