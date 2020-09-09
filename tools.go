package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

// if sth in elements
func FindElement(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// 统一的后端错误返回函数
func errLog(ctx *fasthttp.RequestCtx, serverLog interface{}, clientLog string, httpStatusCode int) {
	log.Println(serverLog)
	ctx.Error(clientLog, httpStatusCode)
}
