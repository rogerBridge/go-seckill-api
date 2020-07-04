package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	//"math/rand"
	"github.com/segmentio/ksuid"
	"strconv"
	"sync"
	"time"
)

// 首先, 初始化redis中待抢购的商品信息
func InitStore() error {
	conn := pool.Get()
	defer conn.Close()

	// 单个用户允许购买的最大数必须小于等于库存数
	if !(limitNumMap["default"] <= storeNum) {
		return errors.New("单个用户允许购买的最大值>库存数, 这是不允许的!")
	}
	// 首先, flush redis
	err := conn.Send("flushdb")
	if err != nil {
		log.Println("flushdb err", err)
		return err
	}
	// 创造store:productId相关数据, 假设: wahaha的商品id是10000
	err = conn.Send("hmset", "store:"+productId, "productName", productName, "productId", productId, "storeNum", storeNum)
	if err != nil {
		log.Println(err, " 创建hash `store:10000`失败")
		return err
	}
	// 创造store:10001 相关的数据
	err = conn.Send("hmset", "store:10001", "productName", "cola", "productId", "10001", "storeNum", storeNum)
	if err != nil {
		log.Println(err, " 创建hash `store:10001`失败")
		return err
	}
	return nil
}

type User struct {
	UserId string `json:"userId"`
}

// 首先查找 productId && purchaseNum 是否还有足够的库存, 然后在看用户是否满足购买的条件
func (u *User) CanBuyIt(productId string, purchaseNum int) (bool, error) {
	if purchaseNum < 1 || purchaseNum > limitNumMap[productId] {
		return false, errors.New("商品数量不合法或者购买商品数量超出限制!")
	}

	if ok, _ := u.UserFilter(productId, purchaseNum); ok {
		return true, nil
	}
	return false, errors.New("不满足userFilter条件!")
}

// 检查用户是否满足购买某种商品的权限
// 比如说一个用户最多可以购买2个
func (u *User) UserFilter(productId string, purchaseNum int) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	// hget 如果数据库中没有这个hash, 那返回一个空值
	r, err := redis.Int(conn.Do("hget", "user:"+u.UserId+":bought", productId))
	if err != nil {
		return true, nil
	}
	if r>=0 && (r+purchaseNum) <= limitNumMap[productId] {
		return true, nil
	}
	return false, errors.New("购买数量过大或者其他错误!")
}

// ksuid
func orderNumberGenerator() string {
	return ksuid.New().String()
}

// 生成订单
func (u *User) orderGenerator(productId string, purchaseNum int) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	// 我只需要知道当前库存减去purchaseNum是否大于等于0就可
	incrString := strconv.Itoa(purchaseNum)
	value, err := redis.Int(conn.Do("hincrby", "store:"+productId, "storeNum", "-"+incrString))
	if err != nil {
		log.Println(err)
		return "", errors.New("减少库存时出现错误!")
	}
	if value < 0 {
		// 比如说客户想要2件, 这里只有一件, 那就在拒绝客户之后, 把之前减掉的库存再加回来
		err := conn.Send("hincrby", "store:"+productId, "storeNum", incrString)
		if err != nil {
			log.Fatalf("%+v 加库存的时候出现了错误!", u)
		}
		return "", errors.New("库存数量不够客户想要的")
	}
	orderNum := orderNumberGenerator()
	ok, err := redis.String(conn.Do("hmset", "user:"+u.UserId+":order:"+orderNum, "orderNum", orderNum, "userId", u.UserId, "productId", productId, "purchaseNum", purchaseNum, "orderDate", time.Now().Format("2006-01-02 15:04:05"), "status", "process"))
	if err != nil {
		log.Printf("%+v", err)
		return "", err
	}
	if ok == "OK" {
		log.Printf("%+v 购买 %s %d件成功", u, productId, purchaseNum)
		return orderNum, nil
	}
	return "", errors.New("other error")
}

// 把用户产生的订单集合起来, 生成key为: `user:userId:orderNumber` value type: list, value: "order:orderNumber[10]" 类型的数据
func (u *User) orderListAdd(orderNum string) error {
	conn := pool.Get()
	defer conn.Close()
	//
	err := conn.Send("rpush", "user:"+u.UserId+":orderNumList", "order:"+orderNum)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 为了提高效率, 还是使用key为:`user:userId:bought` value type: hash, value: "productId:productNum" 类型的数据吧
func (u *User) Bought(productId string, purchaseNum int) error {
	conn := pool.Get()
	defer conn.Close()
	// 首先看用户的已购买的商品信息里面, 是否存在productId这种货物, 如不存在, 则初始化, 若存在, 则增加
	flag, err := redis.Int(conn.Do("hsetnx", "user:"+u.UserId+":bought", productId, purchaseNum))
	if err != nil {
		log.Println(err)
		return err
	}
	// 如果想要购买的物品已经存在, 那就增加购物车里面的商品的数量
	if flag == 0 {
		err = conn.Send("hincrby", "user:"+u.UserId+":bought", productId, purchaseNum)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// redis接收到订单中心返回给我们的取消的订单, 我们需要恢复库存数和改变 user:[userId]:bought 中特定key对应的value
func (u *User) CancelBuy(orderNum string, m *sync.Mutex) error {
	conn := pool.Get()
	defer conn.Close()
	// 查看订单号是否存在? && 状态是否是process
	m.Lock()
	isOrderExist, err := redis.Int(conn.Do("exists", "user:"+u.UserId+":order:"+orderNum))
	if err != nil {
		log.Printf("%+v 查询user:%s:order:%s 时出错!", u, u.UserId, orderNum)
		m.Unlock()
		return err
	}
	if isOrderExist == 0 {
		log.Printf("%+v 查询user:%s:order:%s 时不存在!", u, u.UserId, orderNum)
		m.Unlock()
		return errors.New("系统中没有找到该订单!")
	}
	status, err := redis.String(conn.Do("hget", "user:"+u.UserId+":order:"+orderNum, "status"))
	if err!=nil {
		log.Printf("%+v 获取用户:%s订单:%s合法性的时候出现了错误!", u, u.UserId, orderNum)
		m.Unlock()
		return err
	}
	if status != "process" {
		log.Printf("%+v 订单%s状态不对", u, orderNum)
		m.Unlock()
		return errors.New("订单状态错误!应该是process")
	}
	// 看用户购买记录hash里是否有这件商品
	// 根据订单号找出来商品的productId, purchaseNum
	productId, err := redis.String(conn.Do("hget", "user:"+u.UserId+":order:"+orderNum, "productId"))
	if err != nil {
		log.Printf("hget user:%s:order:%s productId error", u.UserId, orderNum)
		return err
	}
	purchaseNum, err := redis.Int(conn.Do("hget", "user:"+u.UserId+":order:"+orderNum, "purchaseNum"))
	if err != nil {
		log.Printf("hget user:%s:order:%s purchaseNum error", u.UserId, orderNum)
		return err
	}

	isExist, err := redis.Int(conn.Do("hexists", "user:"+u.UserId+":bought", productId))
	if err != nil {
		log.Printf("%+v 查询user:userId:bought时出错!", u)
		m.Unlock()
		return err
	}
	if isExist == 0 {
		log.Printf("%+v 没有购买过%s", u, productId)
		m.Unlock()
		return errors.New("没有购买过的东东, 不可以取消哦~")
	} else {
		// 人家用户真的购买过...
		existPurchaseNum, err := redis.Int(conn.Do("hget", "user:"+u.UserId+":bought", productId))
		if err != nil {
			log.Printf("%+v 获取已购买商品%s时出现错误! %+v", u, productId, err)
			m.Unlock()
			return err
		}
		if !(existPurchaseNum >= purchaseNum) {
			log.Printf("%+v 已购买数量减去登记的购买数量时出现了错误!", u)
			m.Unlock()
			return errors.New("取消购买的数量不能大于购买的数量!")
		}
		// 给这个订单打个tag status:cancel
		_, err = redis.Int(conn.Do("hset", "user:"+u.UserId+":order:"+orderNum, "status", "cancel"))
		if err!=nil {
			log.Printf("%+v 尝试更改订单: %s 状态时出现错误!", u, orderNum)
		}
		// 返还库存
		incrString := strconv.Itoa(purchaseNum)
		err = conn.Send("hincrby", "store:"+productId, "storeNum", incrString)
		if err != nil {
			log.Printf("%+v 取消订单时出错 @store:productId", u)
			m.Unlock()
			return errors.New(u.UserId + "取消订单时出错 @store:productId")
		}
		// 然后, 改变: user:[userId]:bought 这个hash表里面key对应的value
		err = conn.Send("hincrby", "user:"+u.UserId+":bought", productId, "-"+incrString)
		if err != nil {
			log.Printf("%+v 取消订单时出错 @user:[userId]:bought", u)
			m.Unlock()
			return errors.New(u.UserId + "取消订单时出错 @store:productId")
		}
		m.Unlock()
		return nil
	}
}
