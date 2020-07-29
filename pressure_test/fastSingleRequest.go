package main

import (
	//"encoding/json"
	"github.com/valyala/fasthttp"
	"go_redis/pressure_test/jsonStruct"
	"log"
	"net/http"
	"sync"
	"time"
)

func fastSingleRequest(client fasthttp.Client, userId string, productId string, w *sync.WaitGroup, timeStatistics chan float64) (bool, error){
	// 首先, 构造client

	// 构造request body里面的值
	r := jsonStruct.ReqBuy{
		UserId:      userId,
		ProductId:   productId,
		PurchaseNum: 1,
	}
	reqBody, err := r.MarshalJSON()
	if err!=nil {
		w.Done()
		log.Println(err)
		return false, err
	}
	req := &fasthttp.Request{}
	//req := fasthttp.AcquireRequest()
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(URL)
	req.SetBody(reqBody)
	//req.Header.Set("connection", "close")

	resp := &fasthttp.Response{}
	//resp := fasthttp.AcquireResponse()
	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	err = client.Do(req, resp)
	if err!=nil {
		w.Done()
		log.Println("发送请求时:", err)
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
