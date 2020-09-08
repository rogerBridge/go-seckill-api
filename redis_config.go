package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//// 增加了可靠性, 但是降低了性能
//func newSentinelPool() *redis.Pool {
//	sntnl := &sentinel.Sentinel{
//		//Addrs:      []string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
//		Addrs: []string{"127.0.0.1:26381", "127.0.0.1:26382", "127.0.0.1:26383"},
//		MasterName: "mymaster",
//		Dial: func(addr string) (redis.Conn, error) {
//			//timeout := 500 * time.Millisecond
//			//c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
//			c, err := redis.Dial("tcp", addr)
//			if err != nil {
//				return nil, err
//			}
//			return c, nil
//		},
//	}
//	return &redis.Pool{
//		MaxIdle:     20000,
//		MaxActive:   30000,
//		Wait:        true,
//		IdleTimeout: 300 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			masterAddr, err := sntnl.MasterAddr()
//			if err != nil {
//				return nil, err
//			}
//			c, err := redis.Dial("tcp", masterAddr)
//			if err != nil {
//				return nil, err
//			}
//			return c, nil
//		},
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
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

//定义redis pool
var pool = &redis.Pool{
	MaxIdle:     20000,
	IdleTimeout: 300 * time.Second,
	Dial: func() (conn redis.Conn, err error) {
		networkType := "tcp"
		//host := "127.0.0.1"
		host := "redis"
		masterSocket := host + ":6379"
		//// 如果想要保证webapp在redis挂了的时候, 随着sentinel切换而变更socket
		//c, err := redis.Dial("tcp", "127.0.0.1:6400")
		//if err!=nil {
		//	log.Printf("获取redis配置中心数据有误\n")
		//	return nil, err
		//}
		//masterSocket, err := redis.String(c.Do("get", "masterSocket"))
		//if err!=nil {
		//	log.Printf("获取masterSocket时出错\n")
		//	return nil, err
		//}
		con, err := redis.Dial(networkType, masterSocket) //redis.DialReadTimeout(5*time.Second),
		//redis.DialWriteTimeout(5*time.Second),

		if err != nil {
			return nil, err
		}
		return con, nil
	},
}

// 存储订单信息, 购买信息
var pool1 = &redis.Pool{
	MaxIdle:     20000,
	IdleTimeout: 300 * time.Second,
	Dial: func() (conn redis.Conn, err error) {
		networkType := "tcp"
		//host := "127.0.0.1"
		host := "orderInfoRedis"
		masterSocket := host + ":6379"
		con, err := redis.Dial(networkType, masterSocket)
		if err != nil {
			return nil, err
		}
		return con, nil
	},
}
