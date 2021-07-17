/*
这个包存储了rabbitmq server的配置信息
*/

package common

import (
	"bytes"
	"go-seckill/internal/config"
	"go-seckill/internal/logconf"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Ch作为全局变量, 可以被包外引用, rabbitmqServerName是rabbitmqServer容器在redisStore这个网络中的名称, 其他容器可以根据它的名字找到它
var Ch = GetChannel()

// var rabbitmqServerName = "rabbitmq-server"
// var rabbitmqServerUsername = "root"
// var rabbitmqServerPassword = "12345678"
// var rabbitmqServerPort = "5672"
// var rabbitmqServerPath = "/root_vhost"

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "rabbitmq-common"})

//
type RabbitmqServerConfig struct {
	RabbitmqServerName string
	Username           string
	Password           string
	Port               string
	Path               string
}

// 使用viper从rabbitmq_server_config.json文件中读取键值对
func RabbitmqServerConn() *RabbitmqServerConfig {
	viper.SetConfigType("json")
	err := viper.ReadConfig(bytes.NewBuffer(config.ReadConfig("rabbitmq_server_config.json")))
	//viper.SetConfigFile("./config/rabbitmq_server_config.json")
	//err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("当使用viper读取rabbitmq server的配置时, 出现错误: %v", err)
	}
	rabbitmqServerConfigFromJson := viper.GetStringMapString("rabbitmq-server")
	rabbitmqServerConfig := new(RabbitmqServerConfig)
	rabbitmqServerConfig.RabbitmqServerName = rabbitmqServerConfigFromJson["rabbitmq_server_name"]
	rabbitmqServerConfig.Username = rabbitmqServerConfigFromJson["username"]
	rabbitmqServerConfig.Password = rabbitmqServerConfigFromJson["password"]
	rabbitmqServerConfig.Port = rabbitmqServerConfigFromJson["port"]
	rabbitmqServerConfig.Path = rabbitmqServerConfigFromJson["path"]
	return rabbitmqServerConfig
}
