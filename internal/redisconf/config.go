package redisconf

import (
	"go-seckill/internal/logconf"
	"go-seckill/internal/rabbitmq/common"
	"log"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "redismethods"})

var ch = common.Ch

// 初始化goodsRedisInfo和orderInfoRedis实例
func InitialRedis() error {
	// 搞一些闲置的redis连接
	//var wg sync.WaitGroup
	//for i := 0; i < 5000; i++ {
	//	wg.Add(2)
	//	go newConn(&wg, redis_config.Pool.Get())
	//	go newConn(&wg, redis_config.Pool1.Get())
	//}
	//wg.Wait()
	//log.Println("预热redis链接成功")
	//runtime.GOMAXPROCS(runtime.NumCPU())
	err := InitStore()
	if err != nil {
		log.Println(err)
		return err
	}
	err = LoadGoods()
	if err != nil {
		log.Println(err)
		return err
	}
	// 加载MySQL中的limit到全局变量purchaseLimit中
	err = LoadLimits()
	if err != nil {
		log.Println(err)
		return err
	}
	// // 加载MySQL中的goods到全局变量goodMap中
	// err = GoodMap()
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	return nil
}
