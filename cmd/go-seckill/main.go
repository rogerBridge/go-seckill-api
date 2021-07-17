package main

import (
	"github.com/valyala/fasthttp"
	"go-seckill/internal/mysql/shop_orm"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/router"
	"log"
)

// 执行一些测试性的操作
func test() {

}

func main() {
	test()
	shop_orm.InitialMysql()
	redisconf.InitialRedis()

	r := router.ThisRouter()
	log.Println("Listen on :4000")
	log.Fatalln(fasthttp.ListenAndServe(":4000", r.Handler))
}
