package pressuremaker

import (
	"encoding/json"
	"go-seckill/internal/logconf"
	"io/ioutil"
	"log"
	"sort"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "pressuremaker"})

var (
	ConcurrentNum int
	Host          string
	URL           string
)

type config struct {
	ConcurrentNum int    `json:"concurrentNum"`
	Host          string `json:"host"`
	URL           string `json:"url"`
}

func Start() {
	config := loadConfig()
	ConcurrentNum = config.ConcurrentNum
	Host = config.Host
	URL = Host + config.URL
}

// load config from config.json
func loadConfig() *config {
	fileBytes, err := ioutil.ReadFile("pressuremaker/config.json")
	if err != nil {
		log.Printf("加载配置文件失败: %v\n", err)
		panic(err)
	}
	c := new(config)
	err = json.Unmarshal(fileBytes, c)
	if err != nil {
		log.Printf("read json config from file to struct error \n")
		panic(err)
	}
	return c
}

func PlayTimeStatisticsList(timeStatisticsList []float64) {
	// 请求未完成, 中途夭折的请求的数量
	errorNum := ConcurrentNum - len(timeStatisticsList)
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
	log.Printf("0~1s 内处理的请求数量: %d, 占总体请求数量的%.3f%%\n", durationB0And1, 100*float64(durationB0And1)/float64(ConcurrentNum))
}
