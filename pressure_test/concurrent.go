package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)
type reqBuy struct {
	UserId string `json:"userId"`
	ProductId string `json:"productId"`
	PurchaseNum int `json:"purchaseNum"`
}

var URL = "http://127.0.0.1:4000/buy"
func singleRequest(userId string, w *sync.WaitGroup, badStatus chan int) (bool, error){
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
		w.Done()
		return false, err
	}
	if resp.StatusCode != 200 {
		badStatus <- 1
		//close(badStatus)
	}
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Println(err)
		w.Done()
		return false, err
	}
	log.Println(string(respByte))
	w.Done()
	return true, nil
}