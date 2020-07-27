package main

import (
	"github.com/valyala/fasthttp"
	"log"
	//"net/http"
)

func init() {
	err := InitStore()
	if err != nil {
		log.Println(err)
		return
	}
	//// 搞一些闲置的redis连接
	//var wg sync.WaitGroup
	//for i:=0; i<10000; i++ {
	//	wg.Add(1)
	//	go newConn(&wg, i)
	//}
	//defer wg.Wait()
	//log.Println("预热redis链接成功")
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

	mux := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/buy":
			buy(ctx)
		case "/cancelBuy":
			cancelBuy(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	err := fasthttp.ListenAndServe("0.0.0.0:4000", mux)
	if err != nil {
		log.Fatalln(err)
	}
}
