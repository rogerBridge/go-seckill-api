package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
type reqBuy struct {
	UserId string `json:"userId"`
	ProductId string `json:"productId"`
	PurchaseNum int `json:"purchaseNum"`
}

var URL = "http://127.0.0.1:4000/buy"
func singleRequest(userId string) (bool, error){
	// 首先, 构造client
	client := http.Client{
		Timeout:       30 * time.Second,
	}
	// 构造request body里面的值
	r := reqBuy{
		UserId:      userId,
		ProductId:   "10000",
		PurchaseNum: 1,
	}
	reqBody, err := json.Marshal(r)
	if err!=nil {
		log.Println(err)
		return false, err
	}
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	if err!=nil {
		log.Println(err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	// 开始发送请求
	resp, err := client.Do(req)
	if err!=nil {
		return false, err
	}
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Println(err)
		return false, err
	}
	log.Println(string(respByte))
	return true, nil
}
