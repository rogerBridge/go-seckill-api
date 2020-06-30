package main

import (
	"errors"
	"net/http"
	"sync"
)

var orderGeneratorLock sync.Mutex

// 处理用户要购买某种商品时, 提交的参数: userId, productId, productNum 的参数的处理呀
// 使用application/json的方式
func buy(w http.ResponseWriter, r *http.Request) {
	// 请求方法限定为post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		errorHandle(w, errors.New("请求方法不合法!"), 405)
		return
	}

	buyReqPointer, err := decodeBuyReq(r.Body)
	if err!=nil {
		errorHandle(w, errors.New("reqBody 解析为json时出错!"), 500)
		return
	}
	// 一些数据校验部分, 校验用户id, productId, productNum
	u:= new(User)
	u.UserId = buyReqPointer.UserId
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
	if err!=nil {
		c := CommonResponse{
			Code: 8005,
			Msg:  "商品数量不足或者您不满足购买的条件!",
			Data: nil,
		}
		content, err := commonResp(c)
		if err!=nil {
			errorHandle(w, errors.New(err.Error()), 500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		return
	}
	if ok {
		// 生成订单信息
		//orderNum, err := u.orderGenerator(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		_, err := u.orderGenerator(buyReqPointer.ProductId, buyReqPointer.PurchaseNum, &orderGeneratorLock)
		if err!=nil {
			c := CommonResponse{
				Code: 8002,
				Msg:  "生成订单时发生了错误",
				Data: nil,
			}
			content, err := commonResp(c)
			if err!=nil {
				errorHandle(w, errors.New(err.Error()), 500)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
			return
		}
		//// 给用户的orderList里面添加商品表单
		//err = u.orderListAdd(orderNum)
		//if err!=nil {
		//	c := CommonResponse{
		//		Code: 8003,
		//		Msg:  "向用户的orderList里面添加订单时发生了错误",
		//		Data: nil,
		//	}
		//	content, err := commonResp(c)
		//	if err!=nil {
		//		errorHandle(w, errors.New(err.Error()), 500)
		//	}
		//	w.Header().Set("Content-Type", "application/json")
		//	w.Write(content)
		//	return
		//}
		// 给用户的已经购买的商品hash表里面的值添加数量
		err = u.Bought(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err!=nil {
			c := CommonResponse{
				Code: 8004,
				Msg:  "给用户的已经购买的商品hash表单productId添加数量时发生错误!",
				Data: nil,
			}
			content, err := commonResp(c)
			if err!=nil {
				errorHandle(w, errors.New(err.Error()), 500)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
			return
		}
		//w.Header().Set("application/json", "json")
		c := CommonResponse{
			Code: 8001,
			Msg:  "操作成功!",
			Data: nil,
		}
		content, err := commonResp(c)
		if err!=nil {
			errorHandle(w, errors.New(err.Error()), 500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		return
	}
}