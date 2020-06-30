package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// 首先, 初始化redis中待抢购的商品信息
func InitStore() error {
	conn := pool.Get()
	defer conn.Close()

	// 首先, flush redis
	err := conn.Send("flushdb")
	if err != nil {
		log.Println("flushdb err", err)
		return err
	}
	// 创造store:productId相关数据, 假设: wahaha的商品id是10000
	err = conn.Send("hmset", "store:10000", "productName", "wahaha", "productId", "10000", "storeNum", "2000")
	if err != nil {
		log.Println(err, " 创建hash `store:10000`失败")
		return err
	}
	// 创造store:10001 相关的数据
	err = conn.Send("hmset", "store:10001", "productName", "cola", "productId", "10001", "storeNum", "1000")
	if err!=nil {
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
	//conn := pool.Get()
	//defer conn.Close()
	//// 看指定商品的库存是否还充足?
	//leftNum, err := redis.Int(conn.Do("hget", "store:"+productId, "storeNum"))
	////log.Printf("%T, %v", leftNum, leftNum)
	//if err != nil {
	//	log.Println(err)
	//	return false, err
	//}
	//if ok, _ := u.UserFilter(productId, purchaseNum); leftNum-purchaseNum >= 0 && ok {
	//	log.Printf("%+v could buy it", u)
	//	return true, nil
	//}
	//return false, errors.New("商品数量不足, 或者您不满足UserFilter函数要求!")
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
	// 首先看商品数量是否合法?
	if purchaseNum < 1 || purchaseNum > 2 {
		return false, errors.New("商品数量不合法或者购买商品数量超出限制!")
	}
	// 开始执行事务
	v, err := redis.Int(conn.Do("hexists", "user:"+u.UserId+":bought", productId))
	if err != nil {
		log.Println(err)
		return false, err
	}
	// 如果用户没有购买过, 那直接可以购买
	if v == 0 {
		return true, nil
	} else {
		v, err := redis.Int(conn.Do("hget", "user:"+u.UserId+":bought", productId))
		if err != nil {
			log.Println(err)
			return false, errors.New("hget user:userId:bought 时出现错误!")
		}
		if v > 0 && v < limitNum { // 用户只能购买2件, 就是用户在v=1的时候, 还可以购买一次, 就是只可以购买2件
			return true, nil
		}
	}
	return false, errors.New("购买数量过大或者其他错误!")
}

// 开始购买, 创建订单, hash的key名称格式是: order:[randomlen10], 并且将key作为用户orderList 这个list里面的值
func orderNumberGenerator(length int) string {
	// 生成随机数必备
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 生成订单
func (u *User) orderGenerator(productId string, purchaseNum int, m *sync.Mutex) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	// 首先, 查一下库存
	m.Lock()
	leftNum, err := redis.Int(conn.Do("hget", "store:"+productId, "storeNum"))
	if err!=nil {
		log.Println(err)
		m.Unlock()
		return "", errors.New("查询库存时返回语句出现错误!")
	}
	if leftNum <= 0 {
		log.Println("查询到的库存数量不足啊!")
		m.Unlock()
		return "", errors.New("查询到的库存数量不足啊!")
	}
	// 注意啦, 把库存先搞掉 :)
	incrString := "-" + strconv.Itoa(purchaseNum)
	value, err := redis.Int(conn.Do("hincrby", "store:"+productId, "storeNum", incrString))
	m.Unlock()
	if err!=nil {
		log.Println(err)
		return "", errors.New("减少库存时出现错误!")
	}
	if value <0 {
		return "", errors.New("库存告急!")
	}
	// 生成订单信息 key为: `order:[orderId]`, value为:
	//    UserId      string
	//    ProductId   string
	//    OrderNum       int
	//    OrderTime   string
	orderNum := orderNumberGenerator(orderNumLength)
	ok, err := redis.String(conn.Do("hmset", "order:"+orderNum, "userId", u.UserId, "productId", productId, "purchaseNum", purchaseNum, "orderDate", time.Now().Format("2006-01-02 15:04:05")))
	if ok == "OK" {
		log.Printf("%+v 购买 %s %d件成功", u, productId, purchaseNum)
		return orderNum, nil
	}
	if err != nil {
		log.Println(err)
		return "", err
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