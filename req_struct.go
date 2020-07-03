package main

import (
	"encoding/json"
	"io"
)

type BuyReq struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
}

type CancelBuyReq struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
	OrderNum    string `json:"orderNum"`
}

func decodeBuyReq(buyReq io.ReadCloser) (*BuyReq, error) {
	b := new(BuyReq)
	err := json.NewDecoder(buyReq).Decode(b)
	if err != nil {
		return new(BuyReq), err
	}
	return b, nil
}

func decodeCancelBuyReq(cancelBuyReq io.ReadCloser) (*CancelBuyReq, error) {
	b := new(CancelBuyReq)
	err := json.NewDecoder(cancelBuyReq).Decode(b)
	if err != nil {
		return new(CancelBuyReq), err
	}
	return b, nil
}
