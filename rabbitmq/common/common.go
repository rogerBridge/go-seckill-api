package common

import (
	"log"

	"github.com/streadway/amqp"
)

func Errlog(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}

// Ch作为全局变量, 可以被包外引用, rabbitmqServerName是rabbitmqServer容器在redisStore这个网络中的名称, 其他容器可以根据它的名字找到它
var rabbitmqServerName = "rabbitmqServer"
var Ch = GetChannel(rabbitmqServerName)

// 作为客户端, 和rabbitmq server建立信道, 声明exchange, 声明queue, 绑定queue
func GetChannel(rabbitmqServerName string) *amqp.Channel {
	// 注意, rabbitmqServer的初始化, rabbitmqServer是docker容器在network之中的名称
	conn, err := amqp.Dial("amqp://root:12345678@" + rabbitmqServerName + ":5672/root_vhost")
	Errlog(err, "Failed to connect rabbitmqServer")
	//defer conn.Close()

	ch, err := conn.Channel()
	Errlog(err, "Failed to establish channel")
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs", // 这个exchange负责将订单消息发送给处理写mysql.shop.orders的应用
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	Errlog(err, "declare exchange fail")

	q, err := ch.QueueDeclare(
		"sendToMysql", //queue name
		true,
		false,
		false,
		false,
		nil,
	)
	Errlog(err, "declare queue error happen")

	err = ch.QueueBind(
		q.Name,
		"logsRecord", // receive only key match msg publish
		"logs",
		false,
		nil,
	)
	Errlog(err, "bind exchange to queue error")
	return ch
}
