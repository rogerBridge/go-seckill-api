package jsonStruct

import (
	"encoding/json"
	"log"
)

// 一般情况下, 照着这个来做接口的返回值
type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//
func CommonResp(c CommonResponse) ([]byte, error) {
	v, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}
	return v, nil
}

// 正确生成订单后
type OrderResponse struct {
	UserId      string `json:"userId"`
	PurchaseNum int    `json:"purchaseNum"`
	ProductId   string `json:"productId"`
	OrderNum    string `json:"orderNum"`
}
