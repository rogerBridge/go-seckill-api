package pressuremaker

import (
	"net"
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
}
