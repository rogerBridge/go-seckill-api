package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
)

// InitStore 首先, 初始化redis中待抢购的商品信息
func InitStore() error {
	conn := pool.Get()
	defer conn.Close()

	// 首先, flush redis
	err := conn.Send("flushdb")
	if err != nil {
		log.Println("flushdb err", err)
		return err
	}
	//for i:=0; i<len(pList); i++ {
	//	for j:=0; j<pList[i].StoreNum; j++ {
	//		err = conn.Send("rpush", "store:"+pList[i].ProductId+":have", 0)
	//		if err!=nil {
	//			log.Printf("初始化商品库存失败!\n")
	//			return err
	//		}
	//	}
	//}
	for i := 0; i < len(pList); i++ {
		err = conn.Send("hmset", "store:"+pList[i].productID, "productName", pList[i].ProductName, "productId", pList[i].productID, "storeNum", pList[i].StoreNum)
		if err != nil {
			log.Printf("%+v创建hash `store:%s`失败", err, pList[i].productID)
			return err
		}
	}
	log.Printf("store hash 初始化完成!\n")
	return nil
}

// User is a type to be exported
type User struct {
	userID string
}

// CanBuyIt 首先查找 productId && purchaseNum 是否还有足够的库存, 然后在看用户是否满足购买的条件
func (u *User) CanBuyIt(productID string, purchaseNum int) (bool, error) {
	_, ok := limitNumMap[productID]
	if !ok {
		log.Printf("请求的商品不在限购名单中, 不合法")
		return false, errors.New("请求的商品不在限购名单中, 不合法")
	}
	if purchaseNum < 1 || purchaseNum > limitNumMap[productID] {
		return false, errors.New("商品数量不合法或者购买商品数量超出限制")
	}
	if ok, _ := u.UserFilter(productID, purchaseNum); ok {
		return true, nil
	}
	return false, errors.New("不满足userFilter条件")
}

// UserFilter 检查用户是否满足购买某种商品的权限
func (u *User) UserFilter(productID string, purchaseNum int) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	// hget 用户是否已经购买过了?
	r, err := redis.Int(conn.Do("hget", "user:"+u.userID+":bought", productID))
	// 如果用户没有购买过
	if err == redis.ErrNil {
		return true, nil
	}
	// 如果用户已经购买过, 那么这次购买+之前购买的数量不可以超过总体的限制
	if (r + purchaseNum) <= limitNumMap[productID] {
		return true, nil
	}
	return false, errors.New("购买数量过大或者其他错误")
}

// ksuid generate string
func orderNumberGenerator() string {
	return ksuid.New().String()
}

// 生成订单
func (u *User) orderGenerator(productID string, purchaseNum int) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	//// 只要list rpop之后的值不是nil就可以
	//_, err := redis.Int(conn.Do("rpop", "store:"+productId+":have"))
	//if err == redis.ErrNil {
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
	// 生成订单信息可以使用rabbitmqtt, 将订单信息存储到别的redis上面
	// 生成订单信息
	orderNum := orderNumberGenerator()
	ok, err := redis.String(conn.Do("hmset", "user:"+u.userID+":order:"+orderNum, "orderNum", orderNum, "userID", u.userID, "productId", productID, "purchaseNum", purchaseNum, "orderDate", time.Now().Format("2006-01-02 15:04:05"), "status", "process"))
	if err != nil {
		log.Printf("用户%s生成订单过程中出现了错误: %+v\n", u.userID, err)
		return "", err
	}
	if ok == "OK" {
		log.Printf("%+v 购买 %s %d件成功", u, productID, purchaseNum)
		return orderNum, nil
	}
	return "", errors.New("other error")
}

// Bought 用户成功生成订单信息后, 将已购买这个消息存在于数据库中, 下次还想购买的时候, 就会face限制购买数量的规则哦
func (u *User) Bought(productID string, purchaseNum int) error {
	conn := pool.Get()
	defer conn.Close()
	// 首先看用户的已购买的商品信息里面, 是否存在productId这种货物, 如不存在, 则初始化, 若存在, 则增加
	flag, err := redis.Int(conn.Do("hsetnx", "user:"+u.userID+":bought", productID, purchaseNum))
	if err != nil {
		log.Println("给user:bought这个hash添加用户已购买这个信息的时候, 遇到了错误", err)
		return err
	}
	// 如果想要购买的物品已经存在(之前购买过), setnx没办法添加, 那就再添加一遍吧~
	if flag == 0 {
		err = conn.Send("hincrby", "user:"+u.userID+":bought", productID, purchaseNum)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// CancelBuy redis接收到订单中心返回给我们的取消的订单, 我们需要恢复库存数和改变 user:[userID]:bought 中特定key对应的value
func (u *User) CancelBuy(orderNum string) error {
	conn := pool.Get()
	defer conn.Close()
	// 查看订单号是否存在? && 状态是否是process
	isOrderExist, err := redis.Int(conn.Do("exists", "user:"+u.userID+":order:"+orderNum))
	if err != nil {
		log.Printf("%+v 查询user:%s:order:%s 时出错!", u, u.userID, orderNum)
		return err
	}
	if isOrderExist == 0 {
		log.Printf("%+v 查询user:%s:order:%s 时不存在!", u, u.userID, orderNum)
		return errors.New("系统中没有找到该订单")
	}
	status, err := redis.String(conn.Do("hget", "user:"+u.userID+":order:"+orderNum, "status"))
	if err != nil {
		log.Printf("%+v 获取用户:%s订单:%s合法性的时候出现了错误!", u, u.userID, orderNum)
		return err
	}
	if status != "process" {
		log.Printf("%+v 订单%s状态不对", u, orderNum)
		return errors.New("订单状态错误!只有process状态的订单才可以执行退订单操作")
	}
	// 根据订单号找出来商品的productId, purchaseNum
	productID, err := redis.String(conn.Do("hget", "user:"+u.userID+":order:"+orderNum, "productId"))
	if err != nil {
		log.Printf("hget user:%s:order:%s productId error", u.userID, orderNum)
		return err
	}
	purchaseNum, err := redis.Int(conn.Do("hget", "user:"+u.userID+":order:"+orderNum, "purchaseNum"))
	if err != nil {
		log.Printf("hget user:%s:order:%s purchaseNum error", u.userID, orderNum)
		return err
	}

	isExist, err := redis.Int(conn.Do("hexists", "user:"+u.userID+":bought", productID))
	if err != nil {
		log.Printf("%+v 查询user:userID:bought时出错!", u)
		return err
	}
	if isExist == 0 {
		log.Printf("%+v 没有购买过%s", u, productID)
		return errors.New("没有购买过的东东, 不可以取消哦~")
	}
	// 人家用户真的购买过...
	existPurchaseNum, err := redis.Int(conn.Do("hget", "user:"+u.userID+":bought", productID))
	if err != nil {
		log.Printf("%+v 获取已购买商品%s时出现错误! %+v", u, productID, err)
		return err
	}
	if !(existPurchaseNum >= purchaseNum) {
		log.Printf("%+v 已购买数量减去登记的购买数量时出现了错误!", u)
		return errors.New("取消购买的数量不能大于购买的数量")
	}
	// 给这个订单打个tag status:cancel
	_, err = redis.Int(conn.Do("hset", "user:"+u.userID+":order:"+orderNum, "status", "cancel"))
	if err != nil {
		log.Printf("%+v 尝试更改订单: %s 状态时出现错误!", u, orderNum)
	}
	// 返还库存
	incrString := strconv.Itoa(purchaseNum)
	err = conn.Send("hincrby", "store:"+productID, "storeNum", incrString)
	if err != nil {
		log.Printf("%+v 取消订单时出错 @store:productId", u)
		return errors.New(u.userID + "取消订单时出错 @store:productId")
	}
	// 然后, 改变: user:[userID]:bought 这个hash表里面key对应的value
	err = conn.Send("hincrby", "user:"+u.userID+":bought", productID, "-"+incrString)
	if err != nil {
		log.Printf("%+v 取消订单时出错 @user:[userID]:bought", u)
		return errors.New(u.userID + "取消订单时出错 @store:productId")
	}
	return nil
}
