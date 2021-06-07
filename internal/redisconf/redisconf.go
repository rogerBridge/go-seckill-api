package redisconf

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// 增加了可靠性, 但是降低了性能
//func newSentinelPool() *redis_config.Pool {
//	sntnl := &sentinel.Sentinel{
//Addrs:      []string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
//		Addrs: []string{"127.0.0.1:26381", "127.0.0.1:26382", "127.0.0.1:26383"},
//		MasterName: "mymaster",
//		Dial: func(addr string) (redis_config.Conn, error) {
//timeout := 500 * time.Millisecond
//c, err := redis_config.DialTimeout("tcp", addr, timeout, timeout, timeout)
//			c, err := redis_config.Dial("tcp", addr)
//			if err != nil {
//				return nil, err
//			}
//			return c, nil
//		},
//	}
//	return &redis_config.Pool{
//		MaxIdle:     20000,
//		MaxActive:   30000,
//		Wait:        true,
//		IdleTimeout: 300 * time.Second,
//		Dial: func() (redis_config.Conn, error) {
//			masterAddr, err := sntnl.MasterAddr()
//			if err != nil {
//				return nil, err
//			}
//			c, err := redis_config.Dial("tcp", masterAddr)
//			if err != nil {
//				return nil, err
//			}
//			return c, nil
//		},
//		TestOnBorrow: func(c redis_config.Conn, t time.Time) error {
//			if !sentinel.TestRole(c, "master") {
//				return errors.New("Role check failed")
//			} else {
//				return nil
//			}
//		},
//	}
//}
//
//var pool = newSentinelPool()

// Pool 调用会非常频繁, 因为每个请求都可能会询问Pool
// 实例: goodsInfoRedis的连接池
var Pool = &redis.Pool{
	MaxIdle:     20000,
	IdleTimeout: 60 * time.Second,
	Dial: func() (conn redis.Conn, err error) {
		networkType := "tcp"
		//host := "127.0.0.1"
		host := "redis"
		masterSocket := host + ":6379"
		// 如果想要保证webapp在redis挂了的时候, 随着sentinel切换而变更socket
		//c, err := redis_config.Dial("tcp", "127.0.0.1:6400")
		//if err!=nil {
		//	log.Printf("获取redis配置中心数据有误\n")
		//	return nil, err
		//}
		//masterSocket, err := redis_config.String(c.Do("get", "masterSocket"))
		//if err!=nil {
		//	log.Printf("获取masterSocket时出错\n")
		//	return nil, err
		//}
		conn, err = redis.Dial(networkType, masterSocket) //redis_config.DialReadTimeout(5*time.Second),
		//redis_config.DialWriteTimeout(5*time.Second),

		if err != nil {
			return nil, err
		}
		return conn, nil
	},
}

// // 存储订单信息, 购买信息
// // 实例: orderInfoRedis的连接池
// var Pool1 = &redis.Pool{
// 	MaxIdle:     20000,
// 	IdleTimeout: 300 * time.Second,
// 	Dial: func() (conn redis.Conn, err error) {
// 		networkType := "tcp"
// 		//host := "127.0.0.1"
// 		host := "orderInfoRedis"
// 		masterSocket := host + ":6379"
// 		con, err := redis.Dial(networkType, masterSocket)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return con, nil
// 	},
// }

// // 存储token信息
// // 实例: tokenRedis的连接池
// var Pool2 = &redis.Pool{
// 	MaxIdle:     20000,
// 	IdleTimeout: 300 * time.Second,
// 	Dial: func() (conn redis.Conn, err error) {
// 		networkType := "tcp"
// 		//host := "127.0.0.1"
// 		host := "tokenRedis"
// 		masterSocket := host + ":6379"
// 		con, err := redis.Dial(networkType, masterSocket)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return con, nil
// 	},
// }
