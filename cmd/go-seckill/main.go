package main

import (
	"go-seckill/internal/db/shop_orm"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/router"
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	err := shop_orm.InitialMysql()
	if err != nil {
		log.Fatalln(err)
	}

	err = redisconf.InitialRedis()
	if err != nil {
		log.Fatalln(err)
	}

	r := router.ThisRouter()
	log.Println("Listen on :4000")
	log.Fatalln(fasthttp.ListenAndServe(":4000", r.Handler))
}
