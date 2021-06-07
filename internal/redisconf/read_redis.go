package redisconf

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
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
var Pool = RedisConn("goodsInfoRedis")
var Pool1 = RedisConn("orderInfoRedis")
var Pool2 = RedisConn("tokenRedis")

// redisName 指定redis实例的名称,
func RedisConn(redisName string) *redis.Pool {
	viper.SetConfigFile("./config/redis_config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("当使用viper读取redis配置的时候出现错误, 程序崩溃, error: %v\n", err)
	}
	conn := viper.GetStringMapString(redisName)
	networkType := conn["networktype"]
	host := conn["host"]
	port := conn["port"]
	log.Println("viper读取到的redis配置信息是: ", redisName, networkType, host, port)
	return &redis.Pool{
		IdleTimeout: 300 * time.Second,
		MaxIdle:     30000,
		Dial: func() (conn redis.Conn, err error) {
			networkType := networkType
			host := host
			masterSocket := host + ":" + port
			conn, err = redis.Dial(networkType, masterSocket) //redis_config.DialReadTimeout(5*time.Second),
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
