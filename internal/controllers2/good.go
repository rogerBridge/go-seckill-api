/*
管理good table的接口的集合
*/
package controllers2

import (
	"encoding/json"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop_orm"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"

	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
)

// GoodsList ...
// 从goodsInfoRedis中获取商品清单
func GoodsList(ctx *fasthttp.RequestCtx) {
	redisconn := redisconf.Pool.Get()
	defer redisconn.Close()

	reply, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		logger.Warnf("GoodsList: 获取redis store:* 时错先错误 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "从redis中获取商品信息失败",
			Data: nil,
		})
		return
	}
	type good struct {
		ProductName string `redis:"productName"`
		ProductID   int    `redis:"productId"`
		StoreNum    int    `redis:"storeNum"`
	}
	goodsList := make([]*good, 0)
	for k, v := range reply {
		logger.Infof("GoodsList: current good key is: %v, value is: %v", k, v)
		goodsMap, err := redis.Values(redisconn.Do("hgetall", v))
		if err != nil {
			logger.Warnf("GoodsList: 获取商品key: value时出现错误 %v", err)
			utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "获取商品key: value时出现错误",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		//log.Println(goodsMap)
		g := new(good)
		// goodsMap 还是一个good的map啊, 尽管看起来像是[]interface{}
		err = redis.ScanStruct(goodsMap, g)
		if err != nil {
			logger.Warnf("GoodsList: Redis scanStruct error %v", err)
			utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis scanStruct error",
				Data: nil,
			})
			return
		}
		logger.Infof("GoodsList: After Redis scanStruct %v", g)
		goodsList = append(goodsList, g)
	}
	logger.Infof("GoodsList: 获取商品清单成功 %v", goodsList)
	utils.ResponseWithJson(ctx, fasthttp.StatusOK, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "获取商品清单成功",
		Data: goodsList,
	})
}

// AddGood ...
// 添加单个商品
func CreateGood(ctx *fasthttp.RequestCtx) {
	// 这里必须使用mysql的事务
	// 首先, 从接口中获取good的info
	g := new(shop_orm.Good)
	err := json.Unmarshal(ctx.Request.Body(), g)
	if err != nil {
		logger.Warnf("AddGood: Unmarshal request body error message %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "请求的body解析错误",
			Data: nil,
		})
		return
	}
	logger.Infof("解析后的商品属性是: %+v", g)

	// 查看待添加的商品是否存在
	if g.IfGoodExist() {
		logger.Warnf("CreateGood: 要添加的商品已存在, 不得重复添加")
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "要添加的商品已存在, 不得重复添加",
			Data: nil,
		})
		return
	}
	// 手动开启事务
	tx := mysql.Conn2.Begin()

	err = g.CreateGood(tx)
	if err != nil {
		logger.Warnf("CreateGood: 当添加商品时, 错误发生")
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "添加商品时, 错误发生",
			Data: nil,
		})
		return
	}
	// 把数据写入redis, 如果err!=nil, 呼叫mysql回滚
	redisCoon := redisconf.Pool.Get()
	defer redisCoon.Close()
	_, err = redisCoon.Do("hmset", "store:"+g.ProductID, "productName", g.ProductName, "productId", g.ProductID, "storeNum", g.Inventory)
	if err != nil {
		// 把之前写入到mysql中的good信息也回滚了
		err = tx.Rollback().Error
		if err != nil {
			// logger.Warnf("CreateGood: mysql事务回滚失败 %v", err)
			logger.Fatalf("CreateGood: mysql事务回滚失败 %v", err)
			return
		} else {
			logger.Warnf("CreateGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "写入good数据到redis的过程中出现错误",
				Data: nil,
			})
			return
		}
	}
	// 如果添加进redis的时候没有问题, 那就统一执行mysql事务
	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("CreateGood: commit mysql事务失败 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql transaction commit error",
			Data: nil,
		})
	}
	logger.Infof("CreateGood: commit mysql事务成功")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "添加商品成功",
		Data: nil,
	})
}

func UpdateGood(ctx *fasthttp.RequestCtx) {
	// 首先, 校验格式对不对
	g := new(shop_orm.Good)
	err := json.Unmarshal(ctx.Request.Body(), g)
	if err != nil {
		logger.Warnf("UpdateGood: Unmarshal req body error message %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "request body unmarshal 错误",
			Data: nil,
		})
		return
	}
	// 之后, 查找mysql中是否存在这个商品
	isExist := g.IfGoodExist()
	if !isExist {
		logger.Warnf("ModifyGood: 更新的商品%v不存在", g.ProductID)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "更新的商品不存在",
			Data: nil,
		})
		return
	}
	// 如果存在, 开启mysql事务, 修改mysql和redis
	tx := mysql.Conn2.Begin()
	err = g.UpdateGoodInventory(tx)
	if err != nil {
		logger.Warnf("UpdateGood: 添加mysql事务: UpdateGoods时出现错误 %v", err)
		err := tx.Rollback().Error
		if err != nil {
			logger.Fatalf("UpdateGood: 尝试回滚失败 %v", err)
			return
		} else {
			logger.Warnf("UpdateGood: 尝试回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "mysql事务添加失败, 尝试回滚成功",
				Data: nil,
			})
			return
		}
	}

	redisConn := redisconf.Pool.Get()
	defer redisConn.Close()
	_, err = redisConn.Do("hmset", "store:"+g.ProductID, "productName", g.ProductName, "productId", g.ProductID, "storeNum", g.Inventory)
	if err != nil {
		logger.Warnf("UpdateGood: Redis hmset error %v", err)
		// redis里面增加key时出现错误, 尝试回滚mysql
		err = tx.Rollback().Error
		if err != nil {
			logger.Fatalf("UpdateGood: mysql tx rollback error: %+v", err)
			return
		} else {
			logger.Warnf("ModifyGood: mysql tx rollback success")
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis hmset error",
				Data: nil,
			})
			return
		}
	}

	// redis添加成功后, 开始执行事务
	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("UpdateGood: mysql transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql tx commit error",
			Data: nil,
		})
		return
	}
	logger.Infof("UpdateGood: Update good inventory success")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "update good success",
		Data: nil,
	})
}

// DeleteGood ...
// 删除某个商品信息, 将delete_at column覆盖上时间戳
func DeleteGood(ctx *fasthttp.RequestCtx) {
	g := new(shop_orm.Good)
	err := json.Unmarshal(ctx.Request.Body(), g)
	if err != nil {
		logger.Warnf("DeleteGood: Unmarshal req body error %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "解析参数错误",
			Data: nil,
		})
		return
	}
	// 查看要删除的商品是否存在
	isExist := g.IfGoodExist()
	if !isExist {
		logger.Warnf("DeleteGood: 商品: %v不存在", g.ProductID)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "删除的商品不存在",
			Data: nil,
		})
		return
	}
	// 开启mysql transaction
	tx := mysql.Conn2.Begin()
	err = g.DeleteGood(tx)
	if err != nil {
		logger.Warnf("DeleteGood: 添加mysql事务DeleteGood时出现错误 %v", err)
		err = tx.Rollback().Error
		if err != nil {
			logger.Fatalf("DeleteGood: mysql事务回滚失败 %v", err)
			return
		} else {
			logger.Warnf("DeleteGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "添加mysql事务DeleteGoods时出现错误, 回滚成功",
				Data: nil,
			})
			return
		}
	}

	// 将redis中的商品项目删除
	redisConn := redisconf.Pool.Get()
	defer redisConn.Close()
	_, err = redisConn.Do("del", "store:"+g.ProductID)
	if err != nil {
		logger.Warnf("DeleteGood: Redis delete store:%v error message %v", g.ProductID, err)
		err = tx.Rollback().Error
		if err != nil {
			logger.Fatalf("DeleteGood: mysql事务回滚失败 %v", err)
			return
		} else {
			logger.Warnf("DeleteGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis delete store:product error",
				Data: nil,
			})
			return
		}
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("DeleteGood: mysql tx commit error %v", err)
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql tx commit error",
			Data: nil,
		})
		return
	}
	logger.Infof("DeleteGood: commit mysql deleteGoods tx success")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "删除商品成功",
		Data: nil,
	})
}
