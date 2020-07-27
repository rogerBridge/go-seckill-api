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
	if err!=nil {
		log.Fatalln(err)
	}
}