package main

type BuyReq struct {
	UserId     string `json:"userId"`
	ProductId  string `json:"productId"`
	ProductNum int    `json:"productNum"`
}
