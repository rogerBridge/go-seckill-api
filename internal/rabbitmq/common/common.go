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

// As a rabbitmq client establish rabbitmq server 
// 1. establish channel
// 2. declare exchange 
// 3. declare queue
// 4. bind exchange + queue
func GetChannel() *amqp.Channel {
	rabbitmqServerConfig := RabbitmqServerConn()
	// 注意, rabbitmqServer的初始化, rabbitmqServer是docker容器在network之中的名称
	URL := "amqp://" + rabbitmqServerConfig.Username + ":" + rabbitmqServerConfig.Password + "@" + rabbitmqServerConfig.RabbitmqServerName +
		":" + rabbitmqServerConfig.Port + rabbitmqServerConfig.Path
	logger.Infof("mqtt server url is: %v", URL)

	conn, err := amqp.Dial(URL)
	Errlog(err, "Failed to connect rabbitmqServer")
	//defer conn.Close()

	ch, err := conn.Channel()
	Errlog(err, "Failed to establish channel")
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		"order-process", // 这个exchange负责将订单消息发送给处理写mysql.shop.orders的应用
		"direct",        // exchange 的类型
		true,
		false,
		false,
		false,
		nil,
	)
	Errlog(err, "declare exchange fail")

	q, err := ch.QueueDeclare(
		"order", //queue name
		true,
		false,
		false,
		false,
		nil,
	)
	Errlog(err, "declare queue error happen")

	err = ch.QueueBind(
		q.Name,
		"", // receive only key match msg publish
		"order-process",
		false,
		nil,
	)
	Errlog(err, "bind exchange to queue error")
	return ch
}
