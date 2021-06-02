package utils

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// FindElement if val in slice, return its index, true, if not, return -1, false
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
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Warnf("ResponseWithJson: struct to []byte error happen")
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetBody(response)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(statusCode)
}
