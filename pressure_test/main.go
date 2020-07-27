package main

import (
	"log"
	"sort"
	"strconv"
	"sync"
)

// 同时请求的client数量
var concurrentNum = 10000

// 这个包对已经写成的功能模块进行压力测试
func main() {
	var w sync.WaitGroup
	// 时间统计队列
	timeStatistics := make(chan float64, concurrentNum)

	start := 10000
	end := start + concurrentNum

	//  10000人同时抢购"10000"这件商品
	for i := start; i < end; i++ {
		w.Add(1)
		go fastSingleRequest(strconv.Itoa(i), "10000", &w, timeStatistics)
	}

	//for i:=start; i<15000; i++ {
	//	w.Add(2)
	//	// userId 范围10000~14999的同时抢购"10000"和"10001"
	//	go fastSingleRequest(strconv.Itoa(i), "10000", &w, timeStatistics)
	//	go fastSingleRequest(strconv.Itoa(i+5000), "10001", &w, timeStatistics)
	//}

	w.Wait()
	// 关闭时间统计队列, 开始我们的计算!
	close(timeStatistics)

	//t1 := time.Since(t0).Seconds()
	//log.Printf("服务器角度的每秒事务处理量: %.2f, %d个客户端请求总时间段: %.4fs\n", float64(end-start)/t1, concurrentNum, t1)

	// 把统计到的时间节点放置到一个slice中, 写需要计算的函数方法
	timeStatisticsList := make([]float64, 0, concurrentNum)
	for t := range timeStatistics {
		timeStatisticsList = append(timeStatisticsList, t)
	}
	playTimeStatisticsList(timeStatisticsList)
}

func playTimeStatisticsList(timeStatisticsList []float64) {
	// 请求未完成, 中途夭折的请求的数量
	errorNum := concurrentNum - len(timeStatisticsList)
	log.Printf("客户端总共发送请求: %d个, 客户端角度的没有被服务器处理的请求数量:%d", len(timeStatisticsList), errorNum)

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
	sort.Float64s(timeStatisticsList)
	allTime := 0.0
	for _, v := range timeStatisticsList {
		allTime += v
	}
	log.Printf("最大响应时间: %.4f秒, 最小响应时间: %.4f秒, 平均响应时间: %.4f秒, TPS: %.0f\n", timeStatisticsList[len(timeStatisticsList)-1], timeStatisticsList[0], allTime/float64(len(timeStatisticsList)), float64(len(timeStatisticsList))/timeStatisticsList[len(timeStatisticsList)-1])
}