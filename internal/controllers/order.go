package controllers

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop/goods"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
)

// func errorHandle(w http.ResponseWriter, err error, code int) {
// 	log.Println(err)
// 	http.Error(w, err.Error(), code)
// }

//var cancelBuyLock sync.Mutex

// 处理用户要购买某种商品时, 提交的参数: userId, productId, productNum 的参数的处理呀
// 使用application/json的方式
// func test(w http.ResponseWriter, r *http.Request) {
// }

// Buy ...
// 购买商品的接口
func Buy(ctx *fasthttp.RequestCtx) {
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
	err := buyReqPointer.UnmarshalJSON(ctx.PostBody())
	//err := json.Unmarshal(ctx.PostBody(), buyReqPointer)
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
	u := new(User)
	u.userID = buyReqPointer.UserId
	// 判断productId和productNum是否合法
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
	if err != nil {
		logger.Warnf("Buy: user: %+v CanBuyIt error: %v 您不满足购买商品条件", buyReqPointer.UserId, err)
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
			logger.Warnf("Buy: %v when generate order, error message: %v", buyReqPointer.UserId, err)
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
			logger.Warnf("Buy: 给用户:%v已经购买的商品表单productId添加数量时发生错误: %v", buyReqPointer.UserId, err)
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
		logger.Infof("Buy: %v 购买 %v 操作成功", buyReqPointer.UserId, buyReqPointer.ProductId)
		utils.ResponseWithJson(ctx, fasthttp.StatusOK, easyjsonprocess.CommonResponse{
			Code: 8001,
			Msg:  "操作成功",
			Data: easyjsonprocess.OrderResponse{
				UserId:      buyReqPointer.UserId,
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
	u := new(User)
	u.userID = cancelBuyReqPointer.UserId
	err = u.CancelBuy(cancelBuyReqPointer.OrderNum)
	if err != nil {
		logger.Warnf("CancelBuy: 用户: %s 取消订单: %s 时出现错误", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8006,
			Msg:  fmt.Sprintf("用户: %s 取消订单: %s 时出现错误", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum),
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
	logger.Infof("CancelBuy: 用户: %s 取消订单: %s 成功", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum)
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8007,
		Msg:  fmt.Sprintf("用户: %s 取消订单: %s 成功", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum),
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

// SyncGoodsFromRedis2Mysql ...
// 调用这个函数, 立刻同步(redis中存在的商品(一般情况下, 这个时候mysql中也是存在对应的产品的), redis中的数据同步到mysql), 将redis中已变更的商品数据, 同步到mysql中
// 用途: 更新redis中的商品数据到mysql中
func SyncGoodsFromRedis2Mysql(ctx *fasthttp.RequestCtx) {
	redisconn := redisconf.Pool.Get()
	defer redisconn.Close()
	// 首先, 将redis中存在的商品信息同步到mysql中
	reply, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		logger.Warnf("SyncGoodsFromRedis2Mysql: 同步商品清单从redis到mysql出现错误:%v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "更新redis中现有的商品信息到mysql中出现错误",
			Data: nil,
		})
		//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
		return
	}
	type Goods struct {
		ProductName string `redis:"productName"`
		ProductID   int    `redis:"productId"`
		StoreNum    int    `redis:"storeNum"`
	}
	// make, 先设定一个长度为64的指针
	goodsListRedis := make([]*Goods, 0, 256)
	for _, v := range reply {
		logger.Infof("SyncGoodsFromRedis2Mysql: 商品信息: %v", v)
		goodsMap, err := redis.Values(redisconn.Do("hgetall", v))
		if err != nil {
			logger.Warnf("SyncGoodsFromRedis2Mysql: 从redis中获取商品信息 %v 时出现了错误", v)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "从redis中获取商品信息时出现了错误",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		//log.Println(goodsMap)
		g := new(Goods)
		err = redis.ScanStruct(goodsMap, g)
		if err != nil {
			logger.Warnf("SyncGoodsFromRedis2Mysql: Redis scanStruct error: %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis scanStruct error",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		logger.Info("SyncGoodsFromRedis2Mysql: Redis scanStruct is: ", g)
		goodsListRedis = append(goodsListRedis, g)
	}
	// 开始一个mysql事务
	// 将goodsListRedis中的*Goods添加进mysql
	tx, err := mysql.Conn.Begin()
	if err != nil {
		logger.Warnf("SyncGoodsFromRedis2Mysql: 连接mysql时出现错误: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "",
			Data: nil,
		})
		//ctx.Error(err.Error(), 500)
		return
	}
	// 这里必须使用事务, 不能一条一条的搞
	for _, v := range goodsListRedis {
		_, err := tx.Exec("update goods set product_name=?, inventory=? where product_id=?", v.ProductName, v.StoreNum, v.ProductID)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				logger.Warnf("SyncGoodsFromRedis2Mysql: rollback failed error message: %v", errRollback)
				return
			} else if errRollback == nil {
				logger.Warnf("SyncGoodsFromRedis2Mysql: rollback success, goods info: %v", v)
				return
			}
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "",
				Data: nil,
			})
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Warnf("SyncGoodsFromRedis2Mysql: commit error message: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "提交mysql事务时出现错误",
			Data: nil,
		})
		//ctx.Error(err.Error(), 500)
		return
	}
	logger.Infof("SyncGoodsFromRedis2Mysql: 同步redis信息到mysql成功")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "同步redis信息到mysql成功",
		Data: nil,
	})
	//respJson, err := easyjsonprocess.CommonResp(easyjsonprocess.CommonResponse{
	//	Code: 8001,
	//	Msg:  "处理成功",
	//	Data: nil,
	//})
	//if err != nil {
	//	errLog(ctx, err, err.Error(), 500)
	//	return
	//}
	//ctx.Response.SetStatusCode(200)
	//ctx.Response.SetBody(respJson)
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// SyncGoodsFromMysql2Redis ...
// (mysql中存在 && redis中不存在)的商品数据到redis, 这个接口的用处是: mysql中新添加的商品数据, 需要同步到redis中, 同时保证redis中已存在的商品数据不变
// 用途: Mysql中添加了新的商品数据,把它同步到redis中
func SyncGoodsFromMysql2Redis(ctx *fasthttp.RequestCtx) {
	redisconn := redisconf.Pool.Get()
	defer redisconn.Close()
	// 在现有的MySQL表格中找到所有的商品数据, 比对redis的productList, 如果发现有商品不存在于redis中, 就把它添加进去
	storeList, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		logger.Warnf("SyncGoodsFromMysql2Redis: 获取redis中已经存在的商品信息时出现错误")
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "获取redis中已经存在的商品信息出现错误",
			Data: nil,
		})
		//ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	// 为什么这样做呢? 可以减少memory allocate的时间损耗
	storeIDlist := make([]string, 0, 256) // 分离redis中商品的ID出来, 到单独的store id list
	for _, v := range storeList {
		// 注意, 这里使用v[6:]是因为只想要 `store:` 后面的部分
		storeIDlist = append(storeIDlist, v[6:])
	}
	logger.Infof("SyncGoodsFromMysql2Redis: Redis 中已存在的商品信息是: %v", storeIDlist)
	// 获取mysql中的商品信息
	goodsList, err := goods.QueryGoods()
	if err != nil {
		logger.Warnf("SyncGoodsFromMysql2Redis: 获取mysql中商品信息时出现错误 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  fmt.Sprintf("获取mysql中商品信息时出现错误 %v", err),
			Data: nil,
		})
		return
	}
	// 如果从mysql中获取goodsList顺利, 那就通过比对的方式向redis中添加数据
	for _, v := range goodsList {
		_, ok := utils.FindElement(storeIDlist, strconv.Itoa(v.ProductId))
		if !ok {
			// 给redis中添加相关商品数据
			err = redisconn.Send("hmset", "store:"+strconv.Itoa(v.ProductId), "productName", v.ProductName, "productId", v.ProductId, "storeNum", v.Inventory)
			if err != nil {
				logger.Warnf("SyncGoodsFromMysql2Redis: %+v 创建hash 'store:%d' 失败", err, v.ProductId)
				// 这里有风险, 万一给redis添加信息的时候出现错误, 那就凉凉
				// 比如说: 添加到一半, 出错了, redis又没有回滚机制
				utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
					Code: 8500,
					Msg:  "mysql to redis error",
					Data: nil,
				})
				//ctx.Error("给redis添加更新的产品数据出现错误", 500)
				return
			}
		}
	}
	logger.Infof("SyncGoodsFromMysql2Redis: 将mysql中新的商品数据添加到redis中成功")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "将mysql中新添加的数据缓存到redis中成功",
		Data: nil,
	})
	//respJson, err := easyjsonprocess.CommonResp(easyjsonprocess.CommonResponse{
	//	Code: 8001,
	//	Msg:  "处理成功",
	//	Data: nil,
	//})
	//if err != nil {
	//	errLog(ctx, err, err.Error(), 500)
	//	return
	//}
	//ctx.Response.SetStatusCode(200)
	//ctx.Response.SetBody(respJson)
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

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
		//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
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
	//response := easyjsonprocess.CommonResponse{
	//	Code: 8001,
	//	Msg:  "success",
	//	Data: goodsList,
	//}
	//err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	//if err != nil {
	//	log.Println(err)
	//	ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
	//	return
	//}
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// SyncGoodsLimit ...
// 更新商品限制计划
// 例如, 在更新MySQL的限制购买条件后, 若要将商品购买限制同步到app中, 只需要调用goodsLimit这个接口就可以
func SyncGoodsLimit(ctx *fasthttp.RequestCtx) {
	// 加载limit限制计划
	err := LoadLimit()
	if err != nil {
		logger.Warnf("SyncGoodsLimit: 加载limit变量到全局变量purchaseLimit时出现错误 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "加载mysql中限制购买的数据到全局变量purchaseLimit时出现错误",
			Data: nil,
		})
		return
		//ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	logger.Infof("SyncGoodsLimit: 加载limit变量到全局变量purchaseLimit成功")
	utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "加载mysql中限制购买的数据到全局变量purchaseLimit",
		Data: nil,
	})
	//response := easyjsonprocess.CommonResponse{
	//	Code: 8001,
	//	Msg:  "success",
	//	Data: nil,
	//}
	//err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	//if err != nil {
	//	ctx.Error("internel error", fasthttp.StatusInternalServerError)
	//}
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// AddGood ...
// 添加单个商品
func AddGood(ctx *fasthttp.RequestCtx) {
	// 首先, 从接口中获取good的info
	g := new(structure.Goods)
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
	// 查看待添加的商品是否存在
	isExist, err := goods.IsExist(g.ProductId)
	if err != nil {
		logger.Warnf("AddGood: 查找商品是否存在时出现错误: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "查找商品是否存在时出现错误",
			Data: nil,
		})
		return
	}
	if isExist == 1 {
		logger.Warnf("AddGood: 要添加的商品已存在, 不得重复添加")
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "要添加的商品已存在, 不得重复添加",
			Data: nil,
		})
		return
	}
	// mysql 开启事务
	tx, err := mysql.Conn.Begin()
	if err != nil {
		logger.Warnf("AddGood: 开启mysql事务失败 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "服务器开启mysql transaction失败",
			Data: nil,
		})
		return
	}

	// exec mysql transaction
	err = goods.InsertGoods(tx, g.ProductId, g.ProductName, g.Inventory)
	//_, err = tx.Exec("insert goods (product_id, product_name, inventory) values (?,?,?)", g.ProductId, g.ProductName, g.Inventory)
	if err != nil {
		logger.Warnf("AddGood: 插入商品信息到goods表格时出现错误 %v", err)
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("AddGood: 回滚失败, 立刻中止请求")
			return
		} else if err == nil {
			logger.Warnf("AddGood: 事务添加失败, 回滚成功")
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "mysql事务添加失败, 回滚成功",
				Data: nil,
			})
			return
		}
		return
	}
	logger.Info("AddGood: mysql事务添加成功")

	// write data to redis, if err occur, call mysql tx.rollback
	redisCoon := redisconf.Pool.Get()
	defer redisCoon.Close()
	_, err = redisCoon.Do("hmset", "store:"+strconv.Itoa(g.ProductId), "productName", g.ProductName, "productId", g.ProductId, "storeNum", g.Inventory)
	if err != nil {
		// 把之前写入到mysql中的good信息也回滚了
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("AddGood: mysql事务回滚失败 %v", err)
			return
		} else if err == nil {
			logger.Warnf("AddGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "hmset store:product error",
				Data: nil,
			})
			return
		}
		return
	}
	// 如果添加进redis的时候没有问题, 那就统一执行mysql事务
	err = tx.Commit()
	if err != nil {
		logger.Warnf("AddGood: commit mysql事务失败 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql transaction commit error",
			Data: nil,
		})
	}
	logger.Infof("AddGood: commit mysql事务成功")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "add good succuss",
		Data: nil,
	})
}

// ModifyGood ...
// 修改单个商品信息
func ModifyGood(ctx *fasthttp.RequestCtx) {
	// 首先, 校验格式对不对
	g := new(structure.Goods)
	err := json.Unmarshal(ctx.Request.Body(), g)
	if err != nil {
		logger.Warnf("ModifyGood: Unmarshal req body error message %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "请求的body解析错误",
			Data: nil,
		})
		return
	}
	// 之后, 查找mysql中是否存在这个商品
	isExist, err := goods.IsExist(g.ProductId)
	if err != nil {
		logger.Warnf("ModifyGood: 查找商品是否存在时出现错误: %+v\n", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "查找商品是否存在时出现错误",
			Data: nil,
		})
		return
	}
	if isExist != 1 {
		logger.Warnf("ModifyGood: 更新的商品%v不存在", g.ProductId)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "更新的商品不存在",
			Data: nil,
		})
		return
	}
	// 如果存在, 开启mysql事务, 修改mysql和redis
	tx, err := mysql.Conn.Begin()
	if err != nil {
		logger.Warnf("ModifyGood: mysql事务开启失败 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql transaction start fail",
			Data: nil,
		})
		return
	}

	err = goods.UpdateGoods(tx, g.ProductId, g.ProductName, g.Inventory)
	if err != nil {
		logger.Warnf("ModifyGood: 添加mysql事务: UpdateGoods时出现错误 %v", err)
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("ModifyGood: 尝试回滚失败 %v", err)
			return
		} else if err == nil {
			logger.Warnf("ModifyGood: 尝试回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "mysql事务添加失败, 尝试回滚成功",
				Data: nil,
			})
			return
		}
		return
	}

	redisConn := redisconf.Pool.Get()
	defer redisConn.Close()
	_, err = redisConn.Do("hmset", "store:"+strconv.Itoa(g.ProductId), "productName", g.ProductName, "productId", g.ProductId, "storeNum", g.Inventory)
	if err != nil {
		logger.Warnf("ModifyGood: Redis hmset error %v", err)
		// redis里面增加key时出现错误, 尝试回滚mysql
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("ModifyGood: mysql tx rollback error: %+v", err)
			return
		} else if err == nil {
			logger.Warnf("ModifyGood: mysql tx rollback success")
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis hmset error",
				Data: nil,
			})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Warnf("ModifyGood: mysql transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "mysql tx commit error",
			Data: nil,
		})
	}
	logger.Infof("ModifyGood: Modify good success")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "modify good success",
		Data: nil,
	})
}

// DeleteGood ...
// 删除单个商品信息, mysql: is_delete=1 and redis del store:{product_id}
func DeleteGood(ctx *fasthttp.RequestCtx) {
	g := new(structure.GoodDelete)
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
	isExist, err := goods.IsExist(g.ProductId)
	if err != nil {
		logger.Warnf("DeleteGood: 查看删除的商品: %v是否存在出现错误", g.ProductId)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "查看删除的商品是否存在出现错误",
			Data: nil,
		})
		return
	}
	if isExist != 1 {
		logger.Warnf("DeleteGood: 商品: %v不存在", g.ProductId)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "删除的商品不存在",
			Data: nil,
		})
		return
	}
	// 开启mysql transaction
	tx, err := mysql.Conn.Begin()
	if err != nil {
		logger.Warnf("DeleteGood: 开启mysql事务时失败 %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "开启mysql transation 错误",
			Data: nil,
		})
		return
	}
	err = goods.DeleteGoods(tx, g.ProductId)
	if err != nil {
		logger.Warnf("DeleteGood: 添加mysql事务DeleteGoods时出现错误 %v", err)
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("DeleteGood: mysql事务回滚失败 %v", err)
			return
		} else if err == nil {
			logger.Warnf("DeleteGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "添加mysql事务DeleteGoods时出现错误, 回滚成功",
				Data: nil,
			})
			return
		}
	}

	redisConn := redisconf.Pool.Get()
	defer redisConn.Close()
	_, err = redisConn.Do("del", "store:"+strconv.Itoa(g.ProductId))
	if err != nil {
		logger.Warnf("DeleteGood: Redis delete store:%v error message %v", g.ProductId, err)
		err = tx.Rollback()
		if err != nil {
			logger.Warnf("DeleteGood: mysql事务回滚失败 %v", err)
			return
		} else if err == nil {
			logger.Warnf("DeleteGood: mysql事务回滚成功 %v", err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "redis delete store:product error",
				Data: nil,
			})
			return
		}
	}

	err = tx.Commit()
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
