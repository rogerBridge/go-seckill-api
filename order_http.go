package main

import (
	"errors"
	"net/http"
	"sync"
)

var orderGeneratorLock sync.Mutex

var cancelBuyLock sync.Mutex
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
		errorHandle(w, errors.New("reqBody 解析为struct时出错!"), 500)
		return
	}
	// 一些数据校验部分, 校验用户id, productId, productNum
	u:= new(User)
	u.UserId = buyReqPointer.UserId
	// 判断productId和productNum是否合法
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
				Msg:  "库存数量不足呀~",
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

// redis收到后台的请求, 用户取消了订单, 发给redis用户的: userId, productId, purchaseNum,  redis直接操作用户的: user:[userId]:bought 里面key为productId的, 值重置为0
// 这个接口必须由后台调用, 因为我没有做数据校验
func cancelBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		errorHandle(w, errors.New("请求方式不合法!"), 405)
		return
	}

	// 解析: /cancelBuy接口传过来的三个参数
	cancelBuyReqPointer, err := decodeCancelBuyReq(r.Body)
	if err!=nil {
		errorHandle(w, errors.New("reqBody解析到struct时出错!"), 500)
		return
	}
	u := new(User)
	u.UserId = cancelBuyReqPointer.UserId
	err = u.CancelBuy(cancelBuyReqPointer.ProductId, cancelBuyReqPointer.PurchaseNum, &cancelBuyLock)
	if err!=nil {
		c := CommonResponse{
			Code: 8006,
			Msg:  "取消订单时失败!",
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
	c := CommonResponse{
		Code: 8007,
		Msg:  "取消订单成功!",
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