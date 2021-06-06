package common

import (
	"github.com/streadway/amqp"
)

// 一旦出现任何rabbitmq问题, 立刻中止程序
func Errlog(err error, msg string) {
	if err != nil {
		logger.Fatalf("%v: %v", err, msg)
	}
}

// 作为客户端, 和rabbitmq server建立信道, 声明exchange, 声明queue, 绑定queue
func GetChannel(rabbitmqServerName string) *amqp.Channel {
	// 注意, rabbitmqServer的初始化, rabbitmqServer是docker容器在network之中的名称
	URL := "amqp://" + rabbitmqServerUsername + ":" + rabbitmqServerPassword + "@" + rabbitmqServerName + ":" +
		rabbitmqServerPort + rabbitmqServerPath
	logger.Infof("mqtt server url is: %v", URL)

	conn, err := amqp.Dial(URL)
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
