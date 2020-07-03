package main

import (
	"log"
	"strconv"
	"sync"
	"time"
)

// 同时请求的client数量
var concurrentNum = 20000

// 这个包对已经写成的功能模块进行压力测试
func main() {
	t0 := time.Now()
	var w sync.WaitGroup
	// 时间统计队列
	timeStatistics := make(chan float64, concurrentNum)

	start := 10000
	end := start + concurrentNum
	//  20000人同时抢购"10000"这件商品
	for i := start; i < end; i++ {
		w.Add(1)
		go singleRequest(strconv.Itoa(i), "10000", &w, timeStatistics)
	}
	//for i:=start; i<20000; i++ {
	//	w.Add(2)
	//	// userId 范围10000~19999的抢购"10000", userId 范围20000~29999的抢购"10001"
	//	go singleRequest(strconv.Itoa(i), "10000", &w, timeStatistics)
	//	go singleRequest(strconv.Itoa(i+10000), "10001", &w, timeStatistics)
	//}
	w.Wait()
	// 关闭时间统计队列, 开始我们的计算!
	close(timeStatistics)

	t1 := time.Since(t0).Seconds()
	log.Printf("每秒事务处理量: %.2f, %d个客户端请求总时间段: %.4fs", float64(end-start)/t1, concurrentNum, t1)

	// 把统计到的时间节点放置到一个slice中, 写需要计算的函数方法
	timeStatisticsList := make([]float64, 0, 20000)
	for t := range timeStatistics {
		timeStatisticsList = append(timeStatisticsList, t)
	}
	playTimeStatisticsList(timeStatisticsList, t1)

	//// 服务器直接拒绝的情况, 会出现吗阿哈
	//badNum := 0
	//for v := range badStatus {
	//	badNum += v
	//}
	//log.Printf("服务器没有响应的请求数量:", badNum)
}

func playTimeStatisticsList(timeStatisticsList []float64, allTime float64) {
	// 请求未完成, 中途夭折的请求的数量
	errorNum := concurrentNum - len(timeStatisticsList)
	log.Println("无效请求数量:", errorNum)
	durationB0And1 := 0 // 时间间隔在[0,1)
	durationB1And2 := 0 // 时间间隔在[1,2)
	durationB2And3 := 0 // 时间间隔在[2,3)
	durationB3And4 := 0 // 时间间隔在[3,4)
	durationB4And5 := 0 // 时间间隔在[4,5)
	durationBiggerThan5 := 0 // 时间间隔在[5, +++++)
	
	for i := 0; i < len(timeStatisticsList); i++ {
		x := timeStatisticsList[i]
		switch {
		case x>=0 && x<1:
			durationB0And1 += 1
		case x>=1 && x<2:
			durationB1And2 += 1
		case x>=2 && x<3:
			durationB2And3 += 1
		case x>=3 && x<4:
			durationB3And4 += 1
		case x>=4 && x<5:
			durationB4And5 += 1
		case x>=5:
			durationBiggerThan5 += 1
		}
	}
	log.Println("在0~1秒内服务器就有返回的请求数量是:", durationB0And1)
	log.Println("在1~2秒内服务器就有返回的请求数量是:", durationB1And2)
	log.Println("在2~3秒内服务器就有返回的请求数量是:", durationB2And3)
	log.Println("在3~4秒内服务器就有返回的请求数量是:", durationB3And4)
	log.Println("在4~5秒内服务器就有返回的请求数量是:", durationB4And5)
	log.Printf("在5~%.4f秒内服务器就有返回的请求数量是:%d", allTime, durationBiggerThan5)
}