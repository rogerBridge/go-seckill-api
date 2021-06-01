package pressuremaker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"redisplay/logconf"

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

func init() {
	config := loadConfig()
	ConcurrentNum = config.ConcurrentNum
	Host = config.Host
	URL = Host + config.URL
}

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

// type ReqBuy struct {
// 	UserId      string `json:"userId"`
// 	ProductId   string `json:"productId"`
// 	PurchaseNum int    `json:"purchaseNum"`
// }
