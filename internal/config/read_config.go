package config

import (
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"time"
)

// read config file from http static server
func ReadConfig(path string) []byte {
	// fasthttp client construct
	var client = &fasthttp.Client{
		MaxConnsPerHost: 1024, // 一个fasthttp.Client客户端的最大TCP数量, 一般达不到65535就不会报错
		Dial: func(addr string) (conn net.Conn, err error) {
			//return connLocal, err
			return fasthttp.DialTimeout(addr, 30*time.Second) // tcp 层
		},
		ReadTimeout: 60 * time.Second, // http 应用层, 如果tcp建立起来, 但是服务器不给你回应||回应的时间太久, 难道你要一直耗着吗?  当然是关闭http链接啊
	}
	req := fasthttp.Request{}
	req.Header.SetMethod(fasthttp.MethodGet)
	URI := "http://go-seckill-config:3000/" + path
	req.SetRequestURI(URI)

	resp := fasthttp.Response{}
	err := client.Do(&req, &resp)
	if err != nil {
		log.Fatalln(err)
	}
	response := resp.Body()
	//fmt.Println(string(response))
	return response
}
