package controllers2

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"

	"github.com/valyala/fasthttp"
)

// Buy ...
// 购买商品的接口
func Buy(ctx *fasthttp.RequestCtx) {
	logger.Debugf("验证middleauth中间件的确将token中的信息加到了request header中%+v", ctx.Request.Header.String())
	//// 请求方法限定为post
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Response.Header.Set("Allow", fasthttp.MethodPost)
	//	ctx.Error("request method must be post", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方法不合法!"), 405)
	//	return
	//}

	// 使用了easyjson, 据说可以提高marshal, unmarshal的效率
	buyReqPointer := new(easyjsonprocess.BuyReq)
	// err := buyReqPointer.UnmarshalJSON(ctx.PostBody())
	err := json.NewDecoder(ctx.RequestBodyStream()).Decode(buyReqPointer)
	// err := json.Unmarshal(ctx.PostBody(), buyReqPointer)
	if err != nil {
		logger.Warnf("Buy: Decode buy request error: %v", err)
		utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "无法解析客户端发送的body",
			Data: nil,
		})
		//ctx.Error("decode json body error", 500)
		return
	}

	// 一些数据校验部分, 校验用户id, productId, productNum
	u := new(redisconf.User)
	u.Username = buyReqPointer.Username
	// 判断productId和productNum是否合法
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
	if err != nil {
		logger.Warnf("Buy: user: %+v CanBuyIt error: %v 您不满足购买商品条件", buyReqPointer.Username, err)
		//content, err := c.MarshalJSON()
		//content, err := easyjsonprocess.CommonResp(c)
		//if err != nil {
		//	log.Printf("%v\n", err)
		//	_ = utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, easyjsonprocess.CommonResponse{
		//		Code: 8500,
		//		Msg:  "服务器内部错误: struct > []byte 时出现错误",
		//		Data: nil,
		//	})
		//	return
		//}
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8005,
			Msg:  "您购买的商品数量已达到上限或者缺货: " + err.Error(),
			Data: nil,
		})
		return
	}
	if ok {
		// 生成订单信息
		orderNum, err := u.OrderGenerator(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			logger.Warnf("Buy: %v when generate order, error message: %v", buyReqPointer.Username, err)
			utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
				Code: 8002,
				Msg:  "生成订单过程中出现错误:" + err.Error(),
				Data: nil,
			})
			//content, err := c.MarshalJSON()
			////content, err := easyjsonprocess.CommonResp(c)
			//if err != nil {
			//	log.Println(err)
			//	ctx.Error("store num is not enough", 500)
			//	return
			//	//errorHandle(w, errors.New(err.Error()), 500)
			//}
			//ctx.SetContentType("application/json")
			//ctx.SetBody(content)
			////w.Header().Set("Content-Type", "application/json")
			////w.Write(content)
			return
		}

		// 给用户的已经购买的商品hash表里面的值添加数量
		err = u.Bought(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			logger.Warnf("Buy: 给用户:%v已经购买的商品表单productId添加数量时发生错误: %v", buyReqPointer.Username, err)
			utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
				Code: 8004,
				Msg:  "给用户的已经购买的商品hash表单productId添加数量时发生错误: " + err.Error(),
				Data: nil,
			})
			//content, err := c.MarshalJSON()
			////content, err := easyjsonprocess.CommonResp(c)
			//if err != nil {
			//	//errorHandle(w, errors.New(err.Error()), 500)
			//	log.Println()
			//	ctx.Error("add bought list error", 500)
			//	return
			//}
			//ctx.SetContentType("application/json")
			//ctx.SetBody(content)
			////w.Header().Set("Content-Type", "application/json")
			////w.Write(content)
			return
		}

		//w.Header().Set("application/json", "json")
		logger.Infof("Buy: %v 购买 %v 操作成功", buyReqPointer.Username, buyReqPointer.ProductId)
		utils.ResponseWithJson(ctx, fasthttp.StatusOK, easyjsonprocess.CommonResponse{
			Code: 8001,
			Msg:  "操作成功",
			Data: easyjsonprocess.OrderResponse{
				Username:    buyReqPointer.Username,
				PurchaseNum: buyReqPointer.PurchaseNum,
				ProductId:   buyReqPointer.ProductId,
				OrderNum:    orderNum,
			},
		})
		//content, err := c.MarshalJSON()
		////content, err := easyjsonprocess.CommonResp(c)
		//if err != nil {
		//	log.Println(err)
		//	ctx.Error("json marshal error", 500)
		//	return
		//	//errorHandle(w, errors.New(err.Error()), 500)
		//}
		//ctx.SetContentType("application/json")
		////w.Header().Set("Content-Type", "application/json")
		////w.Write(content)
		//ctx.SetBody(content)
		return
	}
}

// CancelBuy 收到后台的请求, 用户取消了订单, 需要用到的参数有: userId, productId, purchaseNum,  redis直接操作用户的: user:[userId]:bought 里面key为productId的, 赋值为0
func CancelBuy(ctx *fasthttp.RequestCtx) {
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Request.Header.Set("Allow", http.MethodPost)
	//	ctx.Error("request method is not supported", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方式不合法!"), 405)
	//	return
	//}

	// 解析: /cancelBuy接口传过来的四个参数, userId, productId, purchaseNum, orderId
	cancelBuyReqPointer := new(easyjsonprocess.CancelBuyReq)
	err := json.Unmarshal(ctx.Request.Body(), cancelBuyReqPointer)
	if err != nil {
		logger.Warnf("CancelBuy: 解析cancelBuyReqPointer时出现错误 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "解析CancelBuyReq时出现错误",
			Data: nil,
		})
		//ctx.Error("decode request body error", 500)
		return
	}

	//cancelBuyReqPointer, err := decodeCancelBuyReq(r.Body)
	//if err!=nil {
	//	errorHandle(w, errors.New("reqBody解析到struct时出错!"), 500)
	//	return
	//}
	u := new(redisconf.User)
	u.Username = cancelBuyReqPointer.Username
	err = u.CancelBuy(cancelBuyReqPointer.OrderNum)
	if err != nil {
		logger.Warnf("CancelBuy: 用户: %s 取消订单: %s 时出现错误", cancelBuyReqPointer.Username, cancelBuyReqPointer.OrderNum)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8006,
			Msg:  fmt.Sprintf("用户: %s 取消订单: %s 时出现错误", cancelBuyReqPointer.Username, cancelBuyReqPointer.OrderNum),
			Data: nil,
		})
		//c := easyjsonprocess.CommonResponse{
		//	Code: 8006,
		//	Msg:  "取消订单时失败!",
		//	Data: nil,
		//}
		//content, err := easyjsonprocess.CommonResp(c)
		//if err != nil {
		//	log.Println("encode resp body to []byte error", err)
		//	ctx.Error("encode resp body to []byte error", 500)
		//	return
		//}
		//ctx.SetContentType("application/json")
		//ctx.SetBody(content)
		////w.Header().Set("Content-Type", "application/json")
		////w.Write(content)
		return
	}
	logger.Infof("CancelBuy: 用户: %s 取消订单: %s 成功", cancelBuyReqPointer.Username, cancelBuyReqPointer.OrderNum)
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8007,
		Msg:  fmt.Sprintf("用户: %s 取消订单: %s 成功", cancelBuyReqPointer.Username, cancelBuyReqPointer.OrderNum),
		Data: nil,
	})
	//content, err := easyjsonprocess.CommonResp(c)
	//if err != nil {
	//	log.Println(err)
	//	ctx.Error("encode resp body to []byte error", 500)
	//	return
	//	//errorHandle(w, errors.New(err.Error()), 500)
	//}
	//ctx.SetContentType("application/json")
	//ctx.SetBody(content)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(content)
}
