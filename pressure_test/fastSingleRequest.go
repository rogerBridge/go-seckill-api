package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/valyala/fasthttp"
)

func fastSingleRequest(userId string, productId string, w *sync.WaitGroup, timeStatistics chan float64) (bool, error){
	// 首先, 构造client
	client := fasthttp.Client{
		ReadTimeout:                   30*time.Second,
		// 如果不加readtimeout的话, 万一服务器没有正确响应客户端请求, 客户端就会一直保持一个长链接, 直到占用完毕你的tcp连接
	}
	// 构造request body里面的值
	r := reqBuy{
		UserId:      userId,
		ProductId:   productId,
		PurchaseNum: 1,
	}
	reqBody, err := json.Marshal(r)
	if err!=nil {
		w.Done()
		log.Println(err)
		return false, err
	}
	req := &fasthttp.Request{}
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(URL)
	req.SetBody(reqBody)

	resp := &fasthttp.Response{}

	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	err = client.Do(req, resp)
	if err!=nil {
		w.Done()
		log.Println(err)
		return false, err
	}
	if resp.StatusCode() != 200 {
		w.Done()
		return false, err
	}
	t1 := time.Since(t0) // 客户端结束发起请求的时间
	timeStatistics <- t1.Seconds() // 将客户端发起请求的时间发送给timeStatistics

	w.Done()
	return true, nil
}
