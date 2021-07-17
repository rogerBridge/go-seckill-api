/*
使用fasthttp制造http请求
*/
package pressuremaker

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"go-seckill/pressuretest/jsonStruct"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// fasthttp client construct
var FastHttpClient = &fasthttp.Client{
	MaxConnsPerHost: 40960, // 一个fasthttp.Client客户端的最大TCP数量, 一般达不到65535就不会报错
	Dial: func(addr string) (conn net.Conn, err error) {
		//return connLocal, err
		return fasthttp.DialTimeout(addr, 30*time.Second) // tcp 层
	},
	ReadTimeout: 60 * time.Second, // http 应用层, 如果tcp建立起来, 但是服务器不给你回应||回应的时间太久, 难道你要一直耗着吗?  当然是关闭http链接啊
}

type Order struct {
	Token       string
	ProductID   int
	PurchaseNum int
}

func (o *Order) CreateOrder(w *sync.WaitGroup, timeStatistics chan float64, errChan chan error) (bool, error) {
	// 首先, 构造client
	client := FastHttpClient
	var URL = "http://127.0.0.1:4000/api/v0/user/order/buy"

	req := &fasthttp.Request{}
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set("Authorization", o.Token)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(URL)
	// 构造request body里面的值
	r := jsonStruct.ReqBuy{
		ProductId:   o.ProductID,
		PurchaseNum: o.PurchaseNum,
	}
	reqBody, err := json.Marshal(r)
	if err != nil {
		logger.Fatalf("Marshal struct to []byte error message %v", err)
		return false, err
	}
	req.SetBody(reqBody)

	resp := &fasthttp.Response{}
	//resp := fasthttp.AcquireResponse()
	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	err = client.Do(req, resp)
	if err != nil {
		errChan <- fmt.Errorf("client do error %v", err)
		w.Done()
		logger.Warnf("发送请求时: %v", err)
		return false, err
	}
	if resp.StatusCode() != 200 {
		errChan <- fmt.Errorf("server response status code error")
		w.Done()
		return false, err
	}
	t1 := time.Since(t0)           // 客户端结束请求的时间
	timeStatistics <- t1.Seconds() // 将客户端整个请求的时间段发送给timeStatistics
	// 请求结束了
	w.Done()
	return true, nil
}
