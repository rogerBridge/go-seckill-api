package controllers2

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"

	"github.com/valyala/fasthttp"
)

func Buy(ctx *fasthttp.RequestCtx) {
	// 使用了easyjson, 据说可以提高marshal, unmarshal的效率
	order := new(redisconf.Order)
	err := json.Unmarshal(ctx.PostBody(), order)
	if err != nil {
		logger.Warnf("Buy: Decode buy request error: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "解析用户的订单请求时参数出现错误",
			Data: nil,
		})
		return
	}
	order.Username = string(ctx.Request.Header.Peek("username"))
	logger.Infof("从auth处得到的username是: %s", order.Username)
	// 一些数据校验部分, 校验用户id, productId, productNum
	// 判断productId和productNum是否合法
	ok, err := order.CanBuyIt()
	if err != nil {
		logger.Warnf("Buy: user: %+v CanBuyIt error: %v 您不满足购买商品条件", order.Username, err)
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8005,
			Msg:  "你不能购买: " + err.Error(),
			Data: nil,
		})
		return
	}
	if ok {
		// 根据productID从redis中获取商品价格
		
		// 生成订单信息
		err := order.OrderGenerator()
		if err != nil {
			logger.Warnf("用户: %s 生成订单时出现错误: %s", order.Username, err)
			utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
				Code: 8200,
				Msg:  "生成订单过程中出现错误:" + err.Error(),
				Data: nil,
			})
			return
		}

		// 给用户的已经购买的商品hash表里面的值添加数量
		err = order.Bought()
		if err != nil {
			logger.Warnf("用户:%s 添加bought时出现错误: %v", order.Username, err)
			utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
				Code: 8200,
				Msg:  "添加bought键值对时出现错误" + err.Error(),
				Data: nil,
			})
			return
		}

		//w.Header().Set("application/json", "json")
		logger.Infof("用户 %s 购买 %v 操作成功", order.Username, order.ProductID)
		utils.ResponseWithJson(ctx, fasthttp.StatusOK, easyjsonprocess.CommonResponse{
			Code: 8001,
			Msg:  "操作成功",
			Data: order,
		})
		return
	}
}

// CancelBuy 收到后台的请求, 用户取消了订单, 需要用到的参数有: username, orderNumber,  redis直接操作用户的: user:[userId]:bought 里面key为productId的, 赋值为0
func CancelBuy(ctx *fasthttp.RequestCtx) {
	order := new(redisconf.Order)
	//cancelBuyReqPointer := new(easyjsonprocess.CancelBuyReq)
	err := json.Unmarshal(ctx.Request.Body(), order)
	if err != nil {
		logger.Warnf("CancelBuy: 解析用户传来的订单号时出现错误: %s", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "当解析用户传来的订单号时出现错误: " + err.Error(),
			Data: nil,
		})
		return
	}
	order.Username = string(ctx.Request.Header.Peek("username"))
	err = order.CancelBuy()
	if err != nil {
		logger.Warnf("用户: %s 取消订单: %s 时出现错误: %s", order.Username, order.OrderNumber, err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  fmt.Sprintf("用户: %s 取消订单: %s 时出现错误: %s", order.Username, order.OrderNumber, err),
			Data: nil,
		})
		return
	}
	logger.Infof("用户: %s 取消订单: %s 成功", order.Username, order.OrderNumber)
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8200,
		Msg:  fmt.Sprintf("用户: %s 取消订单: %s 成功", order.Username, order.OrderNumber),
		Data: nil,
	})
}
