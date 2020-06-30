package main

import (
	"strconv"
	"time"
)

// 这个包对已经写成的功能模块进行压力测试

func main() {
	for i:=10000; i<11000; i++ {
		go singleRequest(strconv.Itoa(i))
	}
	time.Sleep(5*time.Second)
}

