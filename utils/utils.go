package utils

import (
	"encoding/json"
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

// 统一的json信息返回
func ResponseWithJson(ctx *fasthttp.RequestCtx, statusCode int, payload interface{}) {
	//err := json.NewEncoder(ctx.Response.BodyWriter()).Encode(payload)
	response, err := json.Marshal(payload)
	if err!=nil {
		log.Printf("struct to []byte error happen\n")
	}
	ctx.Response.SetBody(response)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(statusCode)
}
