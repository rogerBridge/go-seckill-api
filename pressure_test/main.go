package main

import (
	"log"
	"strconv"
	"sync"
	"time"
)

var concurrentNum = 20000
// 这个包对已经写成的功能模块进行压力测试
func main() {
	t0 := time.Now()
	var w sync.WaitGroup
	badStatus := make(chan int)
	start := 10000
	end := start + concurrentNum
	for i:=start; i<end; i++ {
		w.Add(1)
		go singleRequest(strconv.Itoa(i), &w, badStatus)
	}
	w.Wait()
	t1 := time.Since(t0).Seconds()
	log.Printf("QPS: %.2f", float64(end-start)/t1)
	badNum := 0
	//for {
	//	v, ok := <-badStatus
	//	if ok == false{
	//		break
	//	}
	//	badNum += v
	//}
	for v := range badStatus {
		badNum += v
	}
	log.Printf("服务器没有响应的请求数量:", badNum)
}

