package redisconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-seckill/internal/mysql/shop/orders"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/mysql/shop_orm"
	"go-seckill/internal/rabbitmq/sender"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
)

type User struct {
	Username string
}

// InitStore, 将mysql中现存的商品添加进redis的goodsInfoRedis实例中
func InitStore() error {
	// goodsInfoRedis
	conn := Pool.Get()
	defer conn.Close()

	// orderInfoRedis
	conn1 := Pool1.Get()
	defer conn1.Close()

	// PING PONG
	err := conn.Send("ping")
	if err != nil {
		logger.Warnf("connect to goodsInfoRedis error message: %v", err)
		return err
	}
	err = conn1.Send("ping")
	if err != nil {
		logger.Warnf("connect to orderInfoRedis error message: %v", err)
		return err
	}
	err = conn.Send("flushdb")
	if err != nil {
		logger.Warnf("flushdb goodsInfoRedis error message: %v", err)
		return err
	}
	err = conn1.Send("flushdb")
	if err != nil {
		logger.Warnf("flushdb orderInfoRedis error message: %v", err)
		return err
	}

	g := &shop_orm.Good{}
	goodsList, err := g.QueryGoods()
	if err != nil {
		logger.Warnf("load goods data from mysql.shop.goods error message: %v", goodsList)
		return err
	}
	for i := 0; i < len(goodsList); i++ {
		err = conn.Send("hmset", "store:"+goodsList[i].ProductName, "name", goodsList[i].ProductName, "category", goodsList[i].ProductCategory, "inventory", goodsList[i].Inventory, "price", goodsList[i].Price)
		if err != nil {
			logger.Warnf("hmset store:%v goodsList error message: %v", goodsList[i], err)
			return err
		}
	}
	logger.Info("load data from mysql.shop.goods to goodsInfoRedis successful ")
	//// 加载limit purchase数据, 比如: 这件商品什么时候可以购买, 一个人可以购买多少件?
	//r, err := purchase_limits.QueryPurchaseLimits()
	//if err!=nil {
	//	log.Println(err)
	//}
	//// limitNum[string]LimitPurchase
	//for _, v := range r {
	//	err = conn.Send("hmset", "limit:"+strconv.Itoa(v.ProductId), "limitNum", v.LimitNum, "startPurchaseTime", v.StartPurchaseDatetime, "endPurchaseTime", v.EndPurchaseDatetime)
	//	if err!=nil {
	//		log.Println(err)
	//		return err
	//	}
	//	log.Println(v.ProductId, v.LimitNum, v.StartPurchaseDatetime, v.EndPurchaseDatetime)
	//}
	return nil
}

// LoadLimit ...
// 加载mysql 中的limit到runtime中
// 全局变量, 存储purchase_limits
var purchaseLimit = make(map[string]*shop_orm.PurchaseLimit)

func LoadLimits() error {
	// 每次修改purchaseLimit变量之前, 先清空这个变量
	for k := range purchaseLimit {
		delete(purchaseLimit, k)
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
		purchaseLimit[v.ProductName] = v
	}
	log.Println("加载后的指针型变量purchaseLimit: ", purchaseLimit)
	for _, v := range purchaseLimit {
		logger.Infof("商品id: %v, 购买限制: %v, 开始购买时间: %v, 结束购买时间: %v", v.ProductName, v.LimitNum, v.StartPurchaseTimeStamp, v.StopPurchaseTimeStamp)
		//err = conn.Send("hmset", "limit:"+strconv.Itoa(v.ProductId), "limitNum", v.LimitNum, "startPurchaseTime", v.StartPurchaseDatetime, "endPurchaseTime", v.EndPurchaseDatetime)
		//if err != nil {
		//	log.Println(err)
		//	return err
		//}
	}
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

// CanBuyIt 首先查找 productId && purchaseNum 是否还有足够的库存, 然后在看用户是否满足购买的条件
func (u *User) CanBuyIt(productID string, purchaseNum int) (bool, error) {
	if _, ok := purchaseLimit[productID]; ok {
		if purchaseNum < 1 || purchaseNum > purchaseLimit[productID].LimitNum {
			return false, fmt.Errorf("用户: %v 购买时商品数量小于1或者超过限制", u)
		}
		if !IfPassTimeCheck(purchaseLimit[productID]) {
			return false, fmt.Errorf("用户: %v 不满足商品购买要求", u)
		}
		// now := time.Now()
		// if purchaseLimit[productID].StartPurchaseTimeStamp == 0 || purchaseLimit[productID].StopPurchaseTimeStamp == 0 {

		// }
		// if now.After(purchaseLimit[productID].StopPurchaseTimeStamp) || now.Before(purchaseLimit[productID].StartPurchaseDatetime) {
		// 	return false, fmt.Errorf("用户: %v 购买时间不符合要求", u)
		// }
		if ok, _ := u.UserFilter(productID, purchaseNum, true); ok {
			return true, nil
		}
		return false, fmt.Errorf("用户: %v 购买数量超出限制", u)
	} else {
		if purchaseNum < 1 {
			return false, fmt.Errorf("用户: %v 购买商品数量小于1", u)
		}
		if ok, _ := u.UserFilter(productID, purchaseNum, false); ok {
			return true, nil
		}
		return false, fmt.Errorf("其他错误")
	}
}

// UserFilter 检查用户是否满足购买某种商品的权限
func (u *User) UserFilter(productID string, purchaseNum int, hasLimit bool) (bool, error) {
	// goodsInfoRedis
	conn := Pool.Get()
	defer conn.Close()
	// orderInfoRedis
	conn1 := Pool1.Get()
	defer conn1.Close()
	// 判断商品库存是否还充足?
	inventory, err := redis.Int(conn.Do("hget", "store:"+productID, "storeNum"))
	if err != nil {
		logger.Warnf("UserFilter: %v获取商品:%s时出现错误\n", u, productID)
		return false, err
	}
	if inventory < 1 {
		return false, nil
	}
	// hget 用户是否已经购买过了?
	r, err := redis.Int(conn1.Do("hget", "user:"+u.Username+":bought", productID))
	// 如果用户没有购买过
	if err == redis.ErrNil {
		return true, nil
	}
	// 如果用户购买过了, 看一看是否存在购买限制
	// 是否存在购买限制
	if hasLimit {
		// 用户: 已购买的数量+想要购买的数量<=限制购买的数量
		if (r + purchaseNum) <= purchaseLimit[productID].LimitNum {
			return true, nil
		}
	} else { // 没有购买限制的话, 那就直接购买咯
		return true, nil
	}
	// 如果用户已经购买过, 那么这次购买+之前购买的数量不可以超过总体的限制
	logger.Warnf("UserFilter: %v购买超过限制或其他错误", u)
	return false, fmt.Errorf("%v购买超过限制或其他错误", u)
}

// ksuid generate string, KSUID将这个作为订单编号
// 订单号生成器的逻辑需要改变
func orderNumberGenerator(u *User) string {
	return fmt.Sprintf("%s-%s-%s", strconv.Itoa(int(time.Now().UnixNano())), u.Username, ksuid.New().String())
}

// 生成订单信息, 首先存入redis缓存, 然后发送到mqtt broker
func (u *User) OrderGenerator(productID string, purchaseNum int) (string, error) {
	conn := Pool.Get()
	defer conn.Close()

	// 存放用户订单信息的redis
	conn1 := Pool1.Get()
	defer conn1.Close()
	//// 只要list rpop之后的值不是nil就可以
	//_, err := redisconf.Int(conn.Do("rpop", "store:"+productId+":have"))
	//if err == redisconf.ErrNil {
	//	return "", errors.New("库存不足")
	//}

	// 我只需要知道当前库存减去purchaseNum是否大于等于0就可
	incrString := strconv.Itoa(purchaseNum)
	value, err := redis.Int(conn.Do("hincrby", "store:"+productID, "storeNum", "-"+incrString))
	if err != nil {
		logger.Warnf("OrderGenerator: %v 减库存过程中出现了错误", u)
		return "", fmt.Errorf("%v减库存过程中出现了错误", u)
	}
	if value < 0 {
		// 比如说客户想要2件, 这里只有一件, 那这波操作之后, 库存就成了-1了, 这是不可接受的, 在拒绝客户之后, 把之前减掉的库存再加回来
		//log.Printf("%s用户在购买过程中, 超卖了, 注意哈", u.userID)
		err := conn.Send("hincrby", "store:"+productID, "storeNum", incrString)
		if err != nil {
			logger.Warnf("OrderGenerator: %v 加库存的过程中出现了错误", u)
		}
		return "", fmt.Errorf("%v 库存数量不够客户想要的", u)
	}
	// 生成订单信息可以使用rabbitmqtt, 将订单信息存储到MySQL上面
	// 生成订单信息, 感觉这里如果用: 时间戳+用户ID 的话, 更好
	// 使用ksuid仅仅是偷懒罢了, 哈哈哈
	orderNum := orderNumberGenerator(u)
	now := time.Now()
	ok, err := redis.String(conn1.Do("hmset", "user:"+u.Username+":order:"+orderNum, "orderNum", orderNum, "userID", u.Username, "productId", productID, "purchaseNum", purchaseNum, "orderDate", now.Format("2006-01-02 15:04:05"), "status", "process"))
	if err != nil {
		logger.Warnf("OrderGenerator: 用户 %v 生成订单过程中出现了错误: %+v\n", u, err)
		return "", err
	}
	if ok == "OK" {
		logger.Warnf("OrderGenerator: %v 购买 %s %d件成功", u, productID, purchaseNum)
	}
	// 开始把生成的信息发送给mqtt exchange
	order := new(structure.Orders)
	order.OrderNum = orderNum
	order.UserId = u.Username
	order.ProductId, err = strconv.Atoi(productID)
	if err != nil {
		logger.Warnf("OrderGenerator: %v productID convert from string to int　%v", u, err)
		return "", err
	}
	order.PurchaseNum = purchaseNum
	order.OrderDatetime = now
	order.Status = "process"
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		logger.Warnf("OrderGenerator: %v json marshal error message: %v", u, err)
		return "", err
	}
	// 将生成的订单信息发送给rabbitmq receiver
	err = sender.Send(jsonBytes, ch)
	if err != nil {
		logger.Warnf("OrderGenerator: %v send msg: %v error message: %v", u, jsonBytes, err)
		return "", err
	}
	return orderNum, nil
}

// Bought 用户成功生成订单信息后, 将已购买这个消息存在于数据库中, 下次还想购买的时候, 就会强制限制购买数量的规则哦
func (u *User) Bought(productID string, purchaseNum int) error {
	conn1 := Pool1.Get()
	defer conn1.Close()
	// 首先看用户的已购买的商品信息里面, 是否存在productId这种货物, 如不存在, 则初始化, 若存在, 则增加
	flag, err := redis.Int(conn1.Do("hsetnx", "user:"+u.Username+":bought", productID, purchaseNum))
	if err != nil {
		logger.Warnf("Bought: add bought info to %v:bought error message: %v", u, err)
		return err
	}
	// 如果想要购买的物品已经存在(之前购买过), setnx没办法添加, 那就再添加一遍吧~
	if flag == 0 {
		err = conn1.Send("hincrby", "user:"+u.Username+":bought", productID, purchaseNum)
		if err != nil {
			logger.Warnf("Bought: when %v buy it, already exist, error message: %v", u, err)
			return err
		}
	}
	return nil
}

// CancelBuy redis接收到订单中心返回给我们的取消的订单, 我们需要恢复库存数和改变 user:[userID]:bought 中特定key对应的value
func (u *User) CancelBuy(orderNum string) error {
	conn := Pool.Get()
	defer conn.Close()

	conn1 := Pool1.Get()
	defer conn1.Close()
	// 查看订单号是否存在? && 状态是否是process
	isOrderExist, err := redis.Int(conn1.Do("exists", "user:"+u.Username+":order:"+orderNum))
	if err != nil {
		logger.Warnf("CancelBuy: %+v 查询user:%s:order:%s 时出错!", u, u.Username, orderNum)
		return err
	}
	if isOrderExist == 0 {
		logger.Warnf("CancelBuy: %+v 查询user:%s:order:%s 时不存在!", u, u.Username, orderNum)
		return fmt.Errorf("%v: 系统中没有找到该订单: %v", u, orderNum)
	}
	status, err := redis.String(conn1.Do("hget", "user:"+u.Username+":order:"+orderNum, "status"))
	if err != nil {
		logger.Warnf("CancelBuy: %v 获取用户:%s订单:%s合法性的时候出现了错误", u, u.Username, orderNum)
		return err
	}
	if status != "process" {
		logger.Warnf("CancelBuy: %+v 订单%s状态不对", u, orderNum)
		return fmt.Errorf("%v 订单状态错误!只有process状态的订单才可以执行退订单操作", u)
	}
	// 根据订单号找出来商品的productId, purchaseNum
	productID, err := redis.String(conn1.Do("hget", "user:"+u.Username+":order:"+orderNum, "productId"))
	if err != nil {
		logger.Warnf("CancelBuy: hget user:%v:order:%v productID error", u, orderNum)
		return err
	}
	purchaseNum, err := redis.Int(conn1.Do("hget", "user:"+u.Username+":order:"+orderNum, "purchaseNum"))
	if err != nil {
		logger.Warnf("CancelBuy: hget user:%s:order:%s purchaseNum error", u.Username, orderNum)
		return err
	}

	isExist, err := redis.Int(conn1.Do("hexists", "user:"+u.Username+":bought", productID))
	if err != nil {
		logger.Warnf("CancelBuy: %v 查询user:%v:bought时出错!", u, u.Username)
		return err
	}
	if isExist == 0 {
		logger.Warnf("CancelBuy: user:%v 没有购买过的东东:%v, 不可以取消哦~", u, productID)
		return fmt.Errorf("user:%v 没有购买过的东东:%v, 不可以取消哦~", u, productID)
	}
	// 如果人家用户真的购买过...
	// 那就赶快处理呀, 嘿嘿
	existPurchaseNum, err := redis.Int(conn1.Do("hget", "user:"+u.Username+":bought", productID))
	if err != nil {
		logger.Warnf("CancelBuy: %+v 获取已购买商品%s时出现错误! %+v", u, productID, err)
		return err
	}
	// 如果顾客购买了1件, 却要取消两件, 那就拒绝
	if !(existPurchaseNum >= purchaseNum) {
		logger.Warnf("CancelBuy: %v 取消购买的数量不能大于购买的数量", u)
		return errors.New("取消购买的数量不能大于购买的数量")
	}
	// 给这个订单打个tag status:cancel
	_, err = redis.Int(conn1.Do("hset", "user:"+u.Username+":order:"+orderNum, "status", "cancel"))
	if err != nil {
		logger.Warnf("CancelBuy: %+v 尝试更改订单: %s 状态时出现错误, 错误原因是: %v\n", u, orderNum, err) // 这里应该将订单状态还原, 并且将错误日志记录在案
		return fmt.Errorf("%+v 尝试更改订单: %s 状态时出现错误, 错误原因是: %v\n", u, orderNum, err)
	}
	// 返还库存
	incrString := strconv.Itoa(purchaseNum)
	err = conn.Send("hincrby", "store:"+productID, "storeNum", incrString)
	if err != nil {
		logger.Warnf("CancelBuy: %v 返还库存 %v 时出错 %v", u, productID, err)
		return fmt.Errorf("%v 返还库存 %v 时出错 %v", u, productID, err)
	}
	// 然后, 改变: user:[userID]:bought 这个hash表里面key对应的value
	err = conn1.Send("hincrby", "user:"+u.Username+":bought", productID, "-"+incrString)
	if err != nil {
		logger.Warnf("CancelBuy: %+v 变更bought表: 'user:%s:bought'时出错%v", u, u.Username, err)
		return fmt.Errorf("%+v 变更bought表: 'user:%s:bought'时出错%v", u, u.Username, err)
	}
	// 最后, 将订单信息同步到mysql中, if订单号不唯一, 那就惨了
	if err := orders.UpdateOrders("cancel", orderNum); err != nil {
		logger.Warnf("CancelBuy: mysql处理更新orders表的时候出错 %v", err)
		return err
	}
	logger.Warnf("用户%v取消订单%s成功", u, orderNum)
	return nil
}
