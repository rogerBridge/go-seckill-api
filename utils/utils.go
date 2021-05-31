package utils

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

// FindElement if sth in elements, return its index, true, if not, return -1, false
func FindElement(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// ResponseWithJson http接口的统一的信息返回
func ResponseWithJson(ctx *fasthttp.RequestCtx, statusCode int, payload interface{}) {
	err := json.NewEncoder(ctx.Response.BodyWriter()).Encode(payload)
	if err != nil {
		log.Println(err)
		return
	}
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("struct to []byte error happen\n")
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetBody(response)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(statusCode)
}
