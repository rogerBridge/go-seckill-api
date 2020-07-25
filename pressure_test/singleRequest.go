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

func singleRequest(userId string, productId string, w *sync.WaitGroup, timeStatistics chan float64) (bool, error){
	// 首先, 构造client
	client := http.Client{
		Timeout:       30 * time.Second,
	}
	// 构造request body里面的值
	r := reqBuy{
		UserId:      userId,
		ProductId:   productId,
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
	t0 := time.Now() // 客户端开始发起请求的时间

	resp, err := client.Do(req)
	if err!=nil {
		w.Done()
		return false, err
	}
	// 服务器无效响应
	if resp.StatusCode != 200 {
		return false, err // 在把时间段返回给服务器的时候就已经return了, timestatistics channel里面收不到数值
	}

	t1 := time.Since(t0) // 客户端结束发起请求的时间
	timeStatistics <- t1.Seconds() // 将客户端发起请求的时间发送给timeStatistics


	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Println(err)
		w.Done()
		return false, err
	}
	w.Done()
	log.Println(string(respByte))
	return true, nil
}