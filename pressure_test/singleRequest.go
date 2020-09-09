package main

import (
	"bytes"
	"encoding/json"
	"go_redis/pressure_test/jsonStruct"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var client1 = &http.Client{
	Transport: &http.Transport{
	},
	Timeout: 30 * time.Second,
}

func singleRequest(client *http.Client, userId string, productId string, w *sync.WaitGroup, timeStatistics chan float64) (bool, error) {
	client = client1
	// 构造request body里面的值
	r := jsonStruct.ReqBuy{
		UserId:      userId,
		ProductId:   productId,
		PurchaseNum: 1,
	}
	reqBody, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return false, err
	}
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	resp, err := client.Do(req)
	if err != nil {
		w.Done()
		return false, err
	}
	// 服务器无效响应
	if resp.StatusCode != 200 {
		w.Done() // 如果不使用的话, 万一程序在此处退出, wait函数将阻塞测试程序的运行
		return false, err
	}
	t1 := time.Since(t0)           // 客户端结束发起请求的时间
	timeStatistics <- t1.Seconds() // 将客户端发起请求的时间发送给timeStatistics

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.Done()
		return false, err
	}
	w.Done()
	//log.Println(string(respByte))
	return true, nil
}
