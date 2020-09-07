package main

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"sync"
)

func init() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	err := InitStore()
	if err != nil {
		log.Println(err)
		return
	}
	// // 搞一些闲置的redis连接
	// var wg sync.WaitGroup
	// for i := 0; i < 10000/2; i++ {
	// 	wg.Add(1)
	// 	go newConn(&wg)
	// }
	// defer wg.Wait()
	// log.Println("预热redis链接成功")
}

// 预热一下客户端, 减少之后的redisPool的链接的内存分配建立连接导致的性能消耗
func newConn(w *sync.WaitGroup) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		log.Fatalln(err)
	}
	w.Done()
}

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/buy", buy)
	//// "/cancelBuy" 这个接口只能由后台来调用
	//mux.HandleFunc("/cancelBuy", cancelBuy)
	//log.Println("Listening on 0.0.0.0:4000")
	//err := http.ListenAndServe("0.0.0.0:4000", mux)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	r := router.New()
	//r.Handle(fasthttp.MethodPost, "/buy", buy)
	r.GET("/goodsList", goodsList)
	r.POST("/syncRedis", syncRedis)
	r.POST("/buy", buy)
	r.POST("/cancelBuy", cancelBuy)
	//mux := func(ctx *fasthttp.RequestCtx) {
	//	switch string(ctx.Path()) {
	//	case "/buy":
	//		buy(ctx)
	//	case "/cancelBuy":
	//		cancelBuy(ctx)
	//	default:
	//		ctx.Error("not found", fasthttp.StatusNotFound)
	//	}
	//}
	log.Println("Listen on :4000")
	err := fasthttp.ListenAndServe(":4000", r.Handler)
	if err != nil {
		log.Fatalln(err)
	}
}
