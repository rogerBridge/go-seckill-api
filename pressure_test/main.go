package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"sync"
)

// 同时请求的client数量
//var concurrentNum = 20000
//var socket = "127.0.0.1:4000"
//var URL = fmt.Sprintf("http://%s/buy", socket)

var TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAyNjY3NzYsInVzZXJuYW1lIjoiZmVucm1lbjMifQ.Q4CkFaMF2C0S6l7Csl26JPdnmfLUuh43l81-FhAX7Hg"
var (
	concurrentNum int
	schema        string
	URL           string
)

type config struct {
	ConcurrentNum int    `json:"concurrentNum"`
	Schema        string `json:"schema"`
	URL           string `json:"URL"`
}

func init() {
	config := loadConfig()
	concurrentNum = config.ConcurrentNum
	schema = config.Schema
	URL = schema + config.URL
}

func loadConfig() *config {
	fileBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf("加载配置文件失败: %v\n", err)
		panic(err)
	}
	c := new(config)
	err = json.Unmarshal(fileBytes, c)
	if err != nil {
		log.Printf("反向json失败\n")
		panic(err)
	}
	return c
}

// 这个包对已经写成的功能模块进行压力测试
func main() {
	var w sync.WaitGroup
	// 时间统计队列
	timeStatistics := make(chan float64, concurrentNum)

	start := 0
	end := start + concurrentNum

	//dialer := &net.Dialer{
	//	LocalAddr: &net.TCPAddr{
	//		IP:   []byte{127, 0, 0, 1},
	//		Port: 5555,
	//	},
	//	Timeout: 30 * time.Second,
	//}
	//connLocal, err := dialer.Dial("tcp", "127.0.0.1:4000")
	//if err != nil {
	//	panic(err)
	//}

	//x 人同时抢购"10001"这件商品
	for i := start; i < end; i++ {
		w.Add(1)
		go fastSingleRequest(client2, strconv.Itoa(i), "10001", &w, timeStatistics)
		//go singleRequest(client1, strconv.Itoa(i), "10001", &w, timeStatistics)
	}

	//for i:=start; i<end-10000; i++ {
	//	w.Add(2)
	//	// userId 范围10000~20000的同时抢购"10000"和"10001"
	//	go fastSingleRequest(client2, strconv.Itoa(i), "10001", &w, timeStatistics)
	//	go fastSingleRequest(client2, strconv.Itoa(i+10000), "10002", &w, timeStatistics)
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
			durationB0And1 += 1
		case x >= 1 && x < 2:
			durationB1And2 += 1
		case x >= 2 && x < 3:
			durationB2And3 += 1
		case x >= 3 && x < 4:
			durationB3And4 += 1
		case x >= 4 && x < 5:
			durationB4And5 += 1
		case x >= 5:
			durationBiggerThan5 += 1
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
	log.Printf("0~1s 内处理的请求数量: %d, 占总体请求数量的%.3f%%\n", durationB0And1, 100*float64(durationB0And1)/float64(concurrentNum))
}