package main

import (
	"encoding/json"
	"io"
)

type BuyReq struct {
	UserId     string `json:"userId"`
	ProductId  string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
}

func decodeBuyReq(buyReq io.ReadCloser) (*BuyReq, error){
	b := new(BuyReq)
	err := json.NewDecoder(buyReq).Decode(b)
	if err!=nil {
		return new(BuyReq), err
	}
	return b, nil
}