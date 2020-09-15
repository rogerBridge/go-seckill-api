package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
	"go_redis/mysql/shop/goods"
	"go_redis/mysql/shop/orders"
	"go_redis/mysql/shop/purchase_limits"
	"go_redis/mysql/shop/structure"
	"go_redis/rabbitmq/common"
	"go_redis/rabbitmq/send"
	"go_redis/redis_config"
	"log"
	"strconv"
	"time"
)

var ch = common.Ch

// InitStore 首先, 初始化redis中待抢购的商品信息
func InitStore() error {
	conn := redis_config.Pool.Get()
	defer conn.Close()

	conn1 := redis_config.Pool1.Get()
	defer conn1.Close()

	// PING PONG
	err := conn.Send("ping")
	if err != nil {
		panic("初始化连接失败: conn fail")
	}
	err = conn1.Send("ping")
	if err != nil {
		panic("初始化连接失败: conn fail")
	}
	// 首先, flushdb redis_config
	err = conn.Send("flushdb")
	if err != nil {
		log.Println("flushdb err", err)
		return err
	}
	err = conn1.Send("flushdb")
	if err != nil {
		log.Println("flushdb err", err)
		return err
	}

	goodsList, err := goods.QueryGoods()
	if err != nil {
		panic("从MySQL数据库加载goods数据失败")
	}
	for i := 0; i < len(goodsList); i++ {
		err = conn.Send("hmset", "store:"+strconv.Itoa(goodsList[i].ProductId), "productName", goodsList[i].ProductName, "productId", goodsList[i].ProductId, "storeNum", goodsList[i].Inventory)
		if err != nil {
			log.Printf("%+v创建hash `store:%s`失败", err, goodsList[i].ProductName)
			return err
		}
	}
	log.Printf("从MySQL数据库中加载goods数据到redis中成功!\n")
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

// 加载limit
func LoadLimit() error {
	conn := redis_config.Pool.Get()
	defer conn.Close()
	for k, _ := range purchaseLimit {
		delete(purchaseLimit, k)
	}
	r, err := purchase_limits.QueryPurchaseLimits()
	if err != nil {
		return err
	}
	for _, v := range r {
		// make map本身就是一个指针型变量
		purchaseLimit[strconv.Itoa(v.ProductId)] = v
	}
	log.Println("加载后的指针型变量purchaseLimit: ", purchaseLimit)
	for _, v := range purchaseLimit {
		log.Println(v.ProductId, v.LimitNum, v.StartPurchaseDatetime, v.EndPurchaseDatetime)
		//err = conn.Send("hmset", "limit:"+strconv.Itoa(v.ProductId), "limitNum", v.LimitNum, "startPurchaseTime", v.StartPurchaseDatetime, "endPurchaseTime", v.EndPurchaseDatetime)
		//if err != nil {
		//	log.Println(err)
		//	return err
		//}
	}
	return nil
}

// 全局变量, 存储purchase_limits
var purchaseLimit = make(map[string]*structure.PurchaseLimits)

// User is a type to be exported
type User struct {
	userID string
}

// CanBuyIt 首先查找 productId && purchaseNum 是否还有足够的库存, 然后在看用户是否满足购买的条件
func (u *User) CanBuyIt(productID string, purchaseNum int) (bool, error) {
	if _, ok := purchaseLimit[productID]; ok {
		if purchaseNum < 1 || purchaseNum > purchaseLimit[productID].LimitNum {
			return false, errors.New("商品数量小于1或者购买商品数量超出限制")
		}
		now := time.Now()
		if now.After(purchaseLimit[productID].EndPurchaseDatetime) || now.Before(purchaseLimit[productID].StartPurchaseDatetime) {
			return false, errors.New("购买时间不符合要求")
		}
		if ok, _ := u.UserFilter(productID, purchaseNum, true); ok {
			return true, nil
		}
		return false, errors.New("购买数量过大或者其他错误")
	} else {
		if purchaseNum < 1 {
			return false, errors.New("商品数量小于1")
		}
		if ok, _ := u.UserFilter(productID, purchaseNum, false); ok {
			return true, nil
		}
		return false, errors.New("其他错误")
	}
}

// UserFilter 检查用户是否满足购买某种商品的权限
func (u *User) UserFilter(productID string, purchaseNum int, hasLimit bool) (bool, error) {
	conn := redis_config.Pool.Get()
	defer conn.Close()
	conn1 := redis_config.Pool1.Get()
	defer conn1.Close()
	// 判断商品库存是否还充足?
	inventory, err := redis.Int(conn.Do("hget", "store:"+productID, "storeNum"))
	if err != nil {
		log.Printf("获取商品:%s时出现错误\n", productID)
		return false, err
	}
	if inventory < 1 {
		return false, nil
	}
	// hget 用户是否已经购买过了?
	r, err := redis.Int(conn1.Do("hget", "user:"+u.userID+":bought", productID))
	// 如果用户没有购买过
	if err == redis.ErrNil {
		return true, nil
	}
	// 如果用户购买过了, 看一看是否存在购买限制
	if hasLimit {
		// 用户想要购买的数量+已购买的数量<=限制购买的数量
		if (r + purchaseNum) <= purchaseLimit[productID].LimitNum {
			return true, nil
		}
	} else {
		return true, nil
	}
	// 如果用户已经购买过, 那么这次购买+之前购买的数量不可以超过总体的限制
	return false, errors.New("其他错误")
}

// ksuid generate string
func orderNumberGenerator() string {
	return ksuid.New().String()
}

// 生成订单
func (u *User) orderGenerator(productID string, purchaseNum int) (string, error) {
	conn := redis_config.Pool.Get()
	defer conn.Close()

	// 存放用户订单信息的redis
	conn1 := redis_config.Pool1.Get()
	defer conn1.Close()
	//// 只要list rpop之后的值不是nil就可以
	//_, err := redis_config.Int(conn.Do("rpop", "store:"+productId+":have"))
	//if err == redis_config.ErrNil {
	//	return "", errors.New("库存不足")
	//}

	// 我只需要知道当前库存减去purchaseNum是否大于等于0就可
	incrString := strconv.Itoa(purchaseNum)
	value, err := redis.Int(conn.Do("hincrby", "store:"+productID, "storeNum", "-"+incrString))
	if err != nil {
		log.Println(err)
		return "", errors.New("减库存过程中出现了错误")
	}
	if value < 0 {
		// 比如说客户想要2件, 这里只有一件, 那这波操作之后, 库存就成了-1了, 这是不可接受的, 在拒绝客户之后, 把之前减掉的库存再加回来
		//log.Printf("%s用户在购买过程中, 超卖了, 注意哈", u.userID)
		err := conn.Send("hincrby", "store:"+productID, "storeNum", incrString)
		if err != nil {
			log.Fatalf("%+v 加库存的时候出现了错误", u)
		}
		return "", errors.New("库存数量不够客户想要的")
	}
	// 生成订单信息可以使用rabbitmqtt, 将订单信息存储到MySQL上面
	// 生成订单信息
	orderNum := orderNumberGenerator()
	now := time.Now()
	ok, err := redis.String(conn1.Do("hmset", "user:"+u.userID+":order:"+orderNum, "orderNum", orderNum, "userID", u.userID, "productId", productID, "purchaseNum", purchaseNum, "orderDate", now.Format("2006-01-02 15:04:05"), "status", "process"))
	if err != nil {
		log.Printf("用户%s生成订单过程中出现了错误: %+v\n", u.userID, err)
		return "", err
	}
	if ok == "OK" {
		log.Printf("%+v 购买 %s %d件成功", u, productID, purchaseNum)
	}
	// 开始把生成的信息发送给mqtt exchange
	order := new(structure.Orders)
	order.OrderNum = orderNum
	order.UserId = u.userID
	order.ProductId, err = strconv.Atoi(productID)
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}
	order.PurchaseNum = purchaseNum
	order.OrderDatetime = now
	order.Status = "process"
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}
	//log.Printf("%s", jsonBytes)
	send.Send(jsonBytes, ch)
	return orderNum, nil
}

// Bought 用户成功生成订单信息后, 将已购买这个消息存在于数据库中, 下次还想购买的时候, 就会face限制购买数量的规则哦
func (u *User) Bought(productID string, purchaseNum int) error {
	conn1 := redis_config.Pool1.Get()
	defer conn1.Close()
	// 首先看用户的已购买的商品信息里面, 是否存在productId这种货物, 如不存在, 则初始化, 若存在, 则增加
	flag, err := redis.Int(conn1.Do("hsetnx", "user:"+u.userID+":bought", productID, purchaseNum))
	if err != nil {
		log.Println("给user:bought这个hash添加用户已购买这个信息的时候, 遇到了错误", err)
		return err
	}
	// 如果想要购买的物品已经存在(之前购买过), setnx没办法添加, 那就再添加一遍吧~
	if flag == 0 {
		err = conn1.Send("hincrby", "user:"+u.userID+":bought", productID, purchaseNum)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// CancelBuy redis接收到订单中心返回给我们的取消的订单, 我们需要恢复库存数和改变 user:[userID]:bought 中特定key对应的value
func (u *User) CancelBuy(orderNum string) error {
	conn := redis_config.Pool.Get()
	defer conn.Close()

	conn1 := redis_config.Pool1.Get()
	defer conn1.Close()
	// 查看订单号是否存在? && 状态是否是process
	isOrderExist, err := redis.Int(conn1.Do("exists", "user:"+u.userID+":order:"+orderNum))
	if err != nil {
		log.Printf("%+v 查询user:%s:order:%s 时出错!", u, u.userID, orderNum)
		return err
	}
	if isOrderExist == 0 {
		log.Printf("%+v 查询user:%s:order:%s 时不存在!", u, u.userID, orderNum)
		return errors.New("系统中没有找到该订单")
	}
	status, err := redis.String(conn1.Do("hget", "user:"+u.userID+":order:"+orderNum, "status"))
	if err != nil {
		log.Printf("%+v 获取用户:%s订单:%s合法性的时候出现了错误!", u, u.userID, orderNum)
		return err
	}
	if status != "process" {
		log.Printf("%+v 订单%s状态不对", u, orderNum)
		return errors.New("订单状态错误!只有process状态的订单才可以执行退订单操作")
	}
	// 根据订单号找出来商品的productId, purchaseNum
	productID, err := redis.String(conn1.Do("hget", "user:"+u.userID+":order:"+orderNum, "productId"))
	if err != nil {
		log.Printf("hget user:%s:order:%s productId error", u.userID, orderNum)
		return err
	}
	purchaseNum, err := redis.Int(conn1.Do("hget", "user:"+u.userID+":order:"+orderNum, "purchaseNum"))
	if err != nil {
		log.Printf("hget user:%s:order:%s purchaseNum error", u.userID, orderNum)
		return err
	}

	isExist, err := redis.Int(conn1.Do("hexists", "user:"+u.userID+":bought", productID))
	if err != nil {
		log.Printf("%+v 查询user:userID:bought时出错!", u)
		return err
	}
	if isExist == 0 {
		log.Printf("%+v 没有购买过%s", u, productID)
		return errors.New("没有购买过的东东, 不可以取消哦~")
	}
	// 如果人家用户真的购买过...
	existPurchaseNum, err := redis.Int(conn1.Do("hget", "user:"+u.userID+":bought", productID))
	if err != nil {
		log.Printf("%+v 获取已购买商品%s时出现错误! %+v", u, productID, err)
		return err
	}
	// 有的人不止购买了一次
	if !(existPurchaseNum >= purchaseNum) {
		log.Printf("%+v 已购买数量减去登记的购买数量时出现了错误!", u)
		return errors.New("取消购买的数量不能大于购买的数量")
	}
	// 给这个订单打个tag status:cancel
	_, err = redis.Int(conn1.Do("hset", "user:"+u.userID+":order:"+orderNum, "status", "cancel"))
	if err != nil {
		log.Printf("%+v 尝试更改订单: %s 状态时出现错误, 错误原因是: %v\n", u, orderNum, err) // 这里应该将订单状态还原, 并且将错误日志记录在案
		return errors.New(err.Error())
	}
	// 返还库存
	incrString := strconv.Itoa(purchaseNum)
	err = conn.Send("hincrby", "store:"+productID, "storeNum", incrString)
	if err != nil {
		log.Printf("%+v 取消订单时出错, product id is: %s, cancel store num is: %s\n", productID, u, incrString)
		return errors.New(u.userID + "取消订单时出错")
	}
	// 然后, 改变: user:[userID]:bought 这个hash表里面key对应的value
	err = conn1.Send("hincrby", "user:"+u.userID+":bought", productID, "-"+incrString)
	if err != nil {
		log.Printf("%+v 变更bought表时出错\n", u)
		return errors.New(u.userID + "变更bought表时出错")
	}
	// 最后, 将订单信息同步到mysql中, if订单号不唯一, 那就惨了
	if err := orders.UpdateOrders("cancel", orderNum); err != nil {
		log.Printf("将订单信息发送到mysql时出现错误 %s", err)
		return err
	}
	log.Printf("用户%s取消订单%s成功\n", u.userID, orderNum)
	return nil
}
