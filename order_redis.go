package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)

// 首先, 初始化redis中待抢购的商品信息
func InitStore() error {
	conn := pool.Get()
	defer conn.Close()

	// 创造store:productId相关数据, 假设: wahaha的商品id是10000
	err := conn.Send("hmset", "store:10000", "productName", "wahaha", "productId", "10000", "storeNum", "200")
	if err!=nil {
		log.Println(err, ":创建hash `store:10000`失败")
		return err
	}
	return nil
}

type User struct {
	UserId string `json:"userId"`
}

// 首先查找productId && productNum是否还有足够的库存, 然后在看用户是否满足购买的条件
func (u *User) CanBuyIt(productId string, productNum int) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	leftNum, err := redis.Int(conn.Do("hget", "store:"+productId, "storeNum"))
	//log.Printf("%T, %v", leftNum, leftNum)
	if err!=nil {
		log.Println(err)
		return false, err
	}
	if leftNum - productNum > 0 && u.UserFilter(productId) == true{
		log.Println("could buy it")
		return true, nil
	}
	return false, errors.New("商品数量不足, 或者您不满足UserFilter函数要求!")
}

// 检查用户是否满足购买某种商品的权限
// 过滤掉重复购买的用户, 其他规则可以自己添加呀 :)
func (u *User) UserFilter(productId string) bool {
	return true
}

