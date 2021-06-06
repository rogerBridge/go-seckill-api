package main

import (
	"go-seckill/internal/controllers"
	"go-seckill/internal/router"
	"log"

	"github.com/valyala/fasthttp"
)

func init() {
	start()
	// 搞一些闲置的redis连接
	//var wg sync.WaitGroup
	//for i := 0; i < 5000; i++ {
	//	wg.Add(2)
	//	go newConn(&wg, redis_config.Pool.Get())
	//	go newConn(&wg, redis_config.Pool1.Get())
	//}
	//wg.Wait()
	//log.Println("预热redis链接成功")
}

func start() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	err := controllers.InitStore()
	if err != nil {
		log.Println(err)
		return
	}
	// 加载MySQL中的limit到全局变量和redis中
	err = controllers.LoadLimit()
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	r := router.ThisRouter()
	log.Println("Listen on :4000")
	log.Fatalln(fasthttp.ListenAndServe(":4000", r.Handler))
}
