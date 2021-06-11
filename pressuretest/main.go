package main

import (
	"go-seckill/pressuretest/pressuremaker"
	"strconv"
	"sync"
)

// 这个包对已经写成的功能模块进行压力测试
// 如果对err信息感兴趣的话, 可以单独写一个分析error信息的函数
func main() {
	pressuremaker.Start()
	token, err := pressuremaker.GetToken()
	if err != nil {
		logger.Fatalf("When get token, error message: %v", err)
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
		go pressuremaker.FastSingleRequest(strconv.Itoa(i), "10001", &w, timeStatistics, token, errChan)
		//go singleRequest(client1, strconv.Itoa(i), "10001", &w, timeStatistics)
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
	// 关闭错误统计channel, 开始我们的计算!
	close(errChan)

	// 把统计到的时间节点放置到一个slice中, 写需要计算的函数方法
	timeStatisticsList := make([]float64, 0, pressuremaker.ConcurrentNum)
	for t := range timeStatistics {
		timeStatisticsList = append(timeStatisticsList, t)
	}
	pressuremaker.PlayTimeStatisticsList(timeStatisticsList)

	// 把统计到的错误信息放置到一个slice中, 写出自己需要的函数方法
	errStatisticsList := make([]error, 0, pressuremaker.ConcurrentNum)
	for e := range errChan {
		errStatisticsList = append(errStatisticsList, e)
	}
	// 正式运行的情况下, 把下面的这个替换为你自己写的错误统计函数
	logger.Println("errChan info: ", errStatisticsList)
}
