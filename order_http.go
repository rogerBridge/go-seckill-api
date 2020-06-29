package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// 处理用户要购买某种商品时, 提交的参数: userId, productId, productNum 的参数的处理呀
// 使用application/json的方式
func buy(w http.ResponseWriter, r *http.Request) {
	// 请求方法限定为post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	buyReqPointer := new(BuyReq)
	err := json.NewDecoder(r.Body).Decode(buyReqPointer)
	if err!=nil {
		log.Println(err)
		return
	}
	log.Println(buyReqPointer)

	// 一些数据校验部分, 校验用户id, productId, productNum
	u:= new(User)
	u.UserId = buyReqPointer.UserId
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.ProductNum)
	if err!=nil {
		errorHandle(w, errors.New(fmt.Sprintf("%+v", buyReqPointer)+"商品不足或者您不满足购买条件!"), 500)
		return
	}
	if ok {
		w.Write([]byte("可以购买呀"))
	}
}
