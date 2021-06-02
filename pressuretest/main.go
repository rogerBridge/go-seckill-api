package main

import (
	"log"
	"os"
	"redisplay/pressuretest/pressuremaker"
	"sort"
	"strconv"
	"sync"
)

// 这个包对已经写成的功能模块进行压力测试
// 如果对err信息感兴趣的话, 可以单独写一个分析error信息的函数
func main() {
	token, err := pressuremaker.GetToken()
	if err != nil {
		logger.Warnf("when generate token, error message %v", err)
		os.Exit(-1)
	}
	var w sync.WaitGroup
	// 时间统计channel
	timeStatistics := make(chan float64, pressuremaker.ConcurrentNum)

	start := 0
	end := start + pressuremaker.ConcurrentNum

	//x 人同时抢购"10001"这件商品
	errChan := make(chan error, pressuremaker.ConcurrentNum)
	for i := start; i < end; i++ {
		w.Add(1)
		// 会将所有的error发送给errChan这个channel, 方便之后统计
		go pressuremaker.FastSingleRequest(strconv.Itoa(i), "10004", &w, timeStatistics, token, errChan)
		//go singleRequest(client1, strconv.Itoa(i), "10001", &w, timeStatistics)
	}
	close(errChan)

	// 遍历errChan之中的错误信息
	for err := range errChan {
		logger.Warnf("%v", err)
	}

	//for i:=start; i<end-10000; i++ {
	//	w.Add(2)
	//	// userId 范围10000~20000的同时抢购"10000"和"10001"
	//	go fastSingleRequest(client2, strconv.Itoa(i), "10001", &w, timeStatistics)
	//	go fastSingleRequest(client2, strconv.Itoa(i+10000), "10002", &w, timeStatistics)
	//}

	w.Wait()
	// 关闭时间统计channel, 开始我们的计算!
	close(timeStatistics)

	//t1 := time.Since(t0).Seconds()
	//log.Printf("服务器角度的每秒事务处理量: %.2f, %d个客户端请求总时间段: %.4fs\n", float64(end-start)/t1, concurrentNum, t1)

	// 把统计到的时间节点放置到一个slice中, 写需要计算的函数方法
	timeStatisticsList := make([]float64, 0, pressuremaker.ConcurrentNum)
	for t := range timeStatistics {
		timeStatisticsList = append(timeStatisticsList, t)
	}
	playTimeStatisticsList(timeStatisticsList)
}

func playTimeStatisticsList(timeStatisticsList []float64) {
	// 请求未完成, 中途夭折的请求的数量
	errorNum := pressuremaker.ConcurrentNum - len(timeStatisticsList)
	log.Printf("客户端总共发送请求: %d个, 客户端角度的没有被服务器处理的请求数量:%d", len(timeStatisticsList), errorNum)

	durationB0And1 := 0      // 时间间隔在[0,1)
	durationB1And2 := 0      // 时间间隔在[1,2)
	durationB2And3 := 0      // 时间间隔在[2,3)
	durationB3And4 := 0      // 时间间隔在[3,4)
	durationB4And5 := 0      // 时间间隔在[4,5)
	durationBiggerThan5 := 0 // 时间间隔在[5, +++++)

	for i := 0; i < len(timeStatisticsList); i++ {
		x := timeStatisticsList[i]
		switch {
		case x >= 0 && x < 1:
			durationB0And1++
		case x >= 1 && x < 2:
			durationB1And2++
		case x >= 2 && x < 3:
			durationB2And3++
		case x >= 3 && x < 4:
			durationB3And4++
		case x >= 4 && x < 5:
			durationB4And5++
		case x >= 5:
			durationBiggerThan5++
		}
	}
	log.Println("在0~1秒内服务器就有返回的请求数量是:", durationB0And1)
	log.Println("在1~2秒内服务器就有返回的请求数量是:", durationB1And2)
	log.Println("在2~3秒内服务器就有返回的请求数量是:", durationB2And3)
	log.Println("在3~4秒内服务器就有返回的请求数量是:", durationB3And4)
	log.Println("在4~5秒内服务器就有返回的请求数量是:", durationB4And5)
	log.Println("大于5秒服务器返回的请求数量是:", durationBiggerThan5)
	sort.Float64s(timeStatisticsList)
	allTime := 0.0
	for _, v := range timeStatisticsList {
		allTime += v
	}
	log.Printf("最大响应时间: %.4fms, 最小响应时间: %.4fms, 平均响应时间: %.4fms, TPS: %.0f\n", 1000*timeStatisticsList[len(timeStatisticsList)-1], 1000*timeStatisticsList[0], 1000*allTime/float64(len(timeStatisticsList)), float64(len(timeStatisticsList))/timeStatisticsList[len(timeStatisticsList)-1])
	log.Printf("0~1s 内处理的请求数量: %d, 占总体请求数量的%.3f%%\n", durationB0And1, 100*float64(durationB0And1)/float64(pressuremaker.ConcurrentNum))
}
