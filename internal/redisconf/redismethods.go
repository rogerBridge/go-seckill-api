package redisconf

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop_orm"
	"go-seckill/internal/rabbitmq/sender"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
)

// InitStore, 将mysql中现存的商品添加进redis的goodsInfoRedis实例中
func InitStore() error {
	// goodRedis
	conn := Pool.Get()
	defer conn.Close()

	// orderRedis
	conn1 := Pool1.Get()
	defer conn1.Close()

	// tokenInfoRedis
	conn2 := Pool2.Get()
	defer conn2.Close()

	// PING PONG
	err := conn.Send("ping")
	if err != nil {
		logger.Fatalf("connect to goodRedis error message: %v", err)
		return err
	}
	err = conn1.Send("ping")
	if err != nil {
		logger.Fatalf("connect to orderRedis error message: %v", err)
		return err
	}
	err = conn2.Send("ping")
	if err != nil {
		logger.Fatalf("connect to tokenInfoRedis error message: %v", err)
		return err
	}
	err = conn.Send("flushdb")
	if err != nil {
		logger.Fatalf("flushdb goodRedis error message: %v", err)
		return err
	}
	err = conn1.Send("flushdb")
	if err != nil {
		logger.Fatalf("flushdb orderRedis error message: %v", err)
		return err
	}
	err = conn2.Send("flushdb")
	if err != nil {
		logger.Fatalf("flushdb tokenInfoRedis error message: %v", err)
		return err
	}
	return nil
}

// LoadLimit ...
// 加载mysql 中的limit到runtime中
// 全局变量, 存储purchase_limits
// product_id 是每件商品的唯一标识符
var purchaseLimitMap = make(map[int]*shop_orm.PurchaseLimit)
var goodMap = GoodMap()

func GoodMap() map[int]*shop_orm.Good {
	// 将mysql中的商品信息加载到redis中
	g := new(shop_orm.Good)
	goodList, err := g.QueryGoods()
	if err != nil {
		logger.Fatalf("load goods data from mysql.shop.goods error message: %v", goodList)
	}
	goodListMap := make(map[int]*shop_orm.Good)
	for _, v := range goodList {
		goodListMap[int((v.ID))] = v
	}
	return goodListMap
}

func LoadLimits() error {
	// 每次修改purchaseLimit变量之前, 先清空这个变量
	for k := range purchaseLimitMap {
		delete(purchaseLimitMap, k)
	}
	p := new(shop_orm.PurchaseLimit)
	r, err := p.QueryPurchaseLimits()
	// r, err := purchase_limits.QueryPurchaseLimits()
	if err != nil {
		logger.Warnf("While query purchaseLimit, error message: %v", err)
		return err
	}
	for _, v := range r {
		// make map本身就是一个指针型变量
		purchaseLimitMap[v.ProductID] = v
	}
	logger.Printf("加载后的指针型变量purchaseLimitMap: %+v", purchaseLimitMap)
	for _, v := range purchaseLimitMap {
		logger.Infof("商品id: %v, 购买限制: %v, 开始购买时间: %v, 结束购买时间: %v", v.ProductID, v.LimitNum, v.StartPurchaseTimeStamp, v.StopPurchaseTimeStamp)
	}
	return nil
}

func LoadGoods() error {
	// goodInfoRedis
	conn := Pool.Get()
	defer conn.Close()
	// 将mysql中的商品信息加载到redis中
	g := new(shop_orm.Good)
	goodList, err := g.QueryGoods()
	if err != nil {
		logger.Warnf("load goods data from mysql.shop.goods error message: %v", goodList)
		return err
	}
	for i := 0; i < len(goodList); i++ {
		err = conn.Send("hmset", "store:"+strconv.Itoa(int(goodList[i].ID)), "productID", goodList[i].ID, "name", goodList[i].ProductName, "category", goodList[i].ProductCategory, "inventory", goodList[i].Inventory, "price", goodList[i].Price)
		if err != nil {
			logger.Warnf("hmset store:%v goodsList error message: %v", goodList[i], err)
			return err
		}
	}
	logger.Info("load data from mysql.shop.goods to goodRedis successful ")
	return nil
}

// 检查是否通过时间段检测
func IfPassTimeCheck(p *shop_orm.PurchaseLimit) bool {
	t := time.Now().UnixNano() / 1e6
	// 开始时间和结束时间任意一项等于0, 说明没有时间段的要求
	if p.StartPurchaseTimeStamp == 0 || p.StopPurchaseTimeStamp == 0 {
		return true
	}
	if t >= int64(p.StartPurchaseTimeStamp) && t <= int64(p.StopPurchaseTimeStamp) {
		return true
	}
	return false
}

type Order shop_orm.Order

// CanBuyIt 首先查找 productId && purchaseNum 是否还有足够的库存, 然后在看用户是否满足购买的条件
func (o *Order) CanBuyIt() (bool, error) {
	if _, ok := purchaseLimitMap[o.ProductID]; ok {
		logger.Debugf("进入CanBuyIt > purchaseLimitMap")
		logger.Debugf("%+v", purchaseLimitMap[o.ProductID])
		if o.PurchaseNum < 1 || o.PurchaseNum > purchaseLimitMap[o.ProductID].LimitNum {
			return false, fmt.Errorf("用户: %s 购买时商品数量小于1或者超过限制", o.Username)
		}
		if !IfPassTimeCheck(purchaseLimitMap[o.ProductID]) {
			return false, fmt.Errorf("用户: %s 不满足商品购买要求", o.Username)
		}
		if ok, err := o.CanBuyIt2(true); ok {
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, fmt.Errorf("用户: %s 购买数量超出限制", o.Username)
	} else {
		if o.PurchaseNum < 1 {
			return false, fmt.Errorf("用户: %s 购买商品数量小于1", o.Username)
		}
		if ok, err := o.CanBuyIt2(false); ok {
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			return false, err
		}
	}
}

// UserFilter 检查用户是否满足购买某种商品的权限
func (o *Order) CanBuyIt2(hasLimit bool) (bool, error) {
	// goodRedis
	conn := Pool.Get()
	defer conn.Close()
	// orderRedis
	conn1 := Pool1.Get()
	defer conn1.Close()
	// 判断商品库存是否还充足?
	inventory, err := redis.Int(conn.Do("hget", "store:"+strconv.Itoa(o.ProductID), "inventory"))
	if err != nil {
		logger.Warnf("UserFilter: %s获取商品:%d时出现错误\n", o.Username, o.ProductID)
		return false, err
	}
	if inventory < 1 {
		return false, nil
	}
	if !hasLimit {
		return true, nil
	} else {
		// hget 用户是否已经购买过了?
		r, err := redis.Int(conn1.Do("hget", "user:"+o.Username+":bought", o.ProductID))
		logger.Debugf("用户: %s 购买商品ID: %d", o.Username, o.ProductID)
		// 如果用户没有购买过
		if err == redis.ErrNil {
			return true, nil
		}
		// 如果用户已经购买过, 那么这次购买+之前购买的数量不可以超过总体的限制
		// 用户: 已购买的数量+想要购买的数量<=限制购买的数量
		if (r + o.PurchaseNum) <= purchaseLimitMap[o.ProductID].LimitNum {
			return true, nil
		} else {
			return false, fmt.Errorf("用户: %s购买数量超出限制", o.Username)
		}
	}
}

// ksuid generate string, KSUID将这个作为订单编号
// 订单号生成器的逻辑需要改变
func (o *Order) orderNumberGenerator() {
	orderNum := ksuid.New().String()
	o.OrderNumber = orderNum
	//return fmt.Sprintf("%s-%s-%s", strconv.Itoa(int(time.Now().UnixNano())), u.Username, ksuid.New().String())
}

// 生成订单信息, 首先存入redis缓存, 然后发送到mqtt broker
// orderNum, username, productID, purchaseNum, orderDate, status
func (o *Order) OrderGenerator() error {
	conn := Pool.Get()
	defer conn.Close()

	// orderRedis
	conn1 := Pool1.Get()
	defer conn1.Close()
	//// 只要list rpop之后的值不是nil就可以
	//_, err := redisconf.Int(conn.Do("rpop", "store:"+productId+":have"))
	//if err == redisconf.ErrNil {
	//	return "", errors.New("库存不足")
	//}

	// 我只需要知道当前库存减去purchaseNum是否大于等于0就可
	incrString := strconv.Itoa(o.PurchaseNum)
	value, err := redis.Int(conn.Do("hincrby", "store:"+strconv.Itoa(o.ProductID), "inventory", "-"+incrString))
	if err != nil {
		logger.Warnf("OrderGenerator: %s 减库存过程中出现了错误", o.Username)
		return fmt.Errorf("%s减库存过程中出现了错误", o.Username)
	}
	if value < 0 {
		// 比如说客户想要2件, 这里只有一件, 那这波操作之后, 库存就成了-1了, 这是不可接受的, 在拒绝客户之后, 把之前减掉的库存再加回来
		//log.Printf("%s用户在购买过程中, 超卖了, 注意哈", u.userID)
		err := conn.Send("hincrby", "store:"+strconv.Itoa(o.ProductID), "inventory", incrString)
		if err != nil {
			logger.Fatalf("OrderGenerator: %s 加库存的过程中出现了错误", o.Username)
		}
		return fmt.Errorf("%s 库存数量不够客户想要的", o.Username)
	}
	// 生成订单信息可以使用rabbitmqtt, 将订单信息存储到MySQL上面
	// 生成订单信息, 感觉这里如果用: 时间戳+用户ID 的话, 更好
	// 使用ksuid仅仅是偷懒罢了, 哈哈哈
	o.orderNumberGenerator()
	now := time.Now()
	o.Status = "process"
	o.Price = goodMap[o.ProductID].Price * o.PurchaseNum
	ok, err := redis.String(conn1.Do("hmset", "user:"+o.Username+":order:"+o.OrderNumber, "orderNum", o.OrderNumber, "username", o.Username, "productID", o.ProductID, "purchaseNum", o.PurchaseNum, "orderDate", now.Format("2006-01-02 15:04:05"), "status", o.Status, "price", o.Price))
	if err != nil {
		logger.Warnf("OrderGenerator: 用户 %s 生成订单过程中出现了错误: %v\n", o.Username, err)
		return err
	}
	if ok == "OK" {
		logger.Infof("OrderGenerator: %s 购买 %d %d件成功", o.Username, o.ProductID, o.PurchaseNum)
	}
	// 开始把生成的信息发送给mqtt exchange

	//logger.Debugln(o.ProductID, goodMap[o.ProductID])
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		logger.Warnf("OrderGenerator: json marshal error message: %v", err)
		return err
	}
	// 将生成的订单信息发送给rabbitmq receiver
	err = sender.Send(jsonBytes, ch)
	logger.Infof("生成的订单信息发送给mqtt exchange: %v", jsonBytes)
	if err != nil {
		logger.Warnf("User: %s send msg: %v error: %v", o.Username, jsonBytes, err)
		return err
	}
	return nil
}

// Bought 用户成功生成订单信息后, 将已购买这个消息存在于redis中, 下次还想购买的时候, 就会强制限制购买数量的规则哦
func (o *Order) Bought() error {
	conn1 := Pool1.Get()
	defer conn1.Close()
	// 首先看用户的已购买的商品信息里面, 是否存在productId这种货物, 如不存在, 则初始化, 若存在, 则增加
	flag, err := redis.Int(conn1.Do("hsetnx", "user:"+o.Username+":bought", o.ProductID, o.PurchaseNum))
	if err != nil {
		logger.Warnf("Bought: add bought info to %s:bought error message: %v", o.Username, err)
		return err
	}
	// 如果想要购买的物品已经存在(之前购买过), setnx没办法添加, 那就再添加一遍吧~
	if flag == 0 {
		err = conn1.Send("hincrby", "user:"+o.Username+":bought", o.ProductID, o.PurchaseNum)
		if err != nil {
			logger.Warnf("Bought: when %s buy it, already exist, error: %v", o.Username, err)
			return err
		}
	}
	return nil
}

// CancelBuy redis接收到订单中心返回给我们的取消的订单, 我们需要恢复库存数和改变 user:[username]:bought 中特定key对应的value
func (o *Order) CancelBuy() error {
	conn := Pool.Get()
	defer conn.Close()

	conn1 := Pool1.Get()
	defer conn1.Close()
	// 查看订单号是否存在? && 状态是否是process
	isOrderExist, err := redis.Int(conn1.Do("exists", "user:"+o.Username+":order:"+o.OrderNumber))
	if err != nil {
		logger.Warnf("CancelBuy: Order%+v 查询user:%s:order:%s 时出错!", o, o.Username, o.OrderNumber)
		return err
	}
	if isOrderExist == 0 {
		logger.Warnf("CancelBuy: Order%+v 查询user:%s:order:%s 时不存在!", o, o.Username, o.OrderNumber)
		return fmt.Errorf("%+v: 系统中没有找到该订单: %v", o, o.OrderNumber)
	}
	status, err := redis.String(conn1.Do("hget", "user:"+o.Username+":order:"+o.OrderNumber, "status"))
	if err != nil {
		logger.Warnf("CancelBuy: Order%+v 获取用户:%s订单:%s合法性的时候出现了错误", o, o.Username, o.OrderNumber)
		return err
	}
	if status != "process" {
		logger.Warnf("CancelBuy: %+v 订单%s状态不对", o, o.OrderNumber)
		return fmt.Errorf("%+v 订单状态错误!只有process状态的订单才可以执行退订单操作", o)
	}
	if status == "process" {
		// 根据订单号找出来商品的productId, purchaseNum
		productID, err := redis.String(conn1.Do("hget", "user:"+o.Username+":order:"+o.OrderNumber, "productID"))
		if err != nil {
			logger.Warnf("CancelBuy: hget user:%v:order:%v productID error", o, o.OrderNumber)
			return err
		}
		purchaseNum, err := redis.Int(conn1.Do("hget", "user:"+o.Username+":order:"+o.OrderNumber, "purchaseNum"))
		if err != nil {
			logger.Warnf("CancelBuy: hget user:%s:order:%s purchaseNum error", o.Username, o.OrderNumber)
			return err
		}
		// 查询用户是否真的购买过这个商品
		isExist, err := redis.Int(conn1.Do("hexists", "user:"+o.Username+":bought", productID))
		if err != nil {
			logger.Warnf("CancelBuy: %+v 查询user:%v:bought时出错!", o, o.Username)
			return err
		}
		if isExist == 0 {
			logger.Warnf("CancelBuy: user:%v 没有购买过的东东:%v, 不可以取消哦~", o.Username, productID)
			return fmt.Errorf("user:%s 没有购买过的东东:%v, 不可以取消哦~", o.Username, productID)
		}
		// 给这个订单打个tag status:cancel
		_, err = redis.Int(conn1.Do("hset", "user:"+o.Username+":order:"+o.OrderNumber, "status", "cancel"))
		if err != nil {
			logger.Warnf("CancelBuy: %+v 尝试更改订单: %s 状态时出现错误, 错误原因是: %v", o, o.OrderNumber, err) // 这里应该将订单状态还原, 并且将错误日志记录在案
			return fmt.Errorf("%+v 尝试更改订单: %s 状态时出现错误, 错误原因是: %v\n", o, o.OrderNumber, err)
		}
		// 打status成功, 开始执行返还库存
		incrString := strconv.Itoa(purchaseNum)
		err = conn.Send("hincrby", "store:"+productID, "inventory", incrString)
		if err != nil {
			logger.Warnf("CancelBuy: %+v 返还库存 %v 时出错 %v", o, productID, err)
			return fmt.Errorf("%+v 返还库存 %v 时出错 %v", o, productID, err)
		}
		// 然后, 改变: user:[userID]:bought 这个hash表里面key对应的value
		err = conn1.Send("hincrby", "user:"+o.Username+":bought", productID, "-"+incrString)
		if err != nil {
			logger.Warnf("CancelBuy: %+v 变更bought表: 'user:%s:bought'时出错%v", o, o.Username, err)
			return fmt.Errorf("%+v 变更bought表: 'user:%s:bought'时出错%v", o, o.Username, err)
		}
		// 最后, 将订单信息同步到mysql中, if订单号不唯一, 那就惨了
		// 之后再加一个订单号是否唯一的校验吧
		// 这里就不加mqtt了, 嘿嘿
		o.Status = "cancel"
		if err := o.UpdateOrderStatus(); err != nil {
			return err
		}
		logger.Warnf("用户%s取消订单%s成功", o.Username, o.OrderNumber)
		return nil
	}
	return fmt.Errorf("未知错误, 哇啦啦")
}

func (o *Order) UpdateOrderStatus() error {
	tx := mysql.Conn2.Begin()
	if err := tx.Model(&Order{}).Where("order_number = ?", o.OrderNumber).Update("status", o.Status).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
