package receive

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go_redis/mysql/shop/orders"
	"go_redis/mysql/shop/structure"
	"go_redis/rabbitmq/common"
	"log"
)

func Receive() {
	conn, err := amqp.Dial("amqp://root:12345678@my_rabbit:5672/root_vhost")
	common.Errlog(err, "Failed to connect my_rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	common.Errlog(err, "Failed to establish channel")
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	common.Errlog(err, "Channel Qos error ")

	err = ch.ExchangeDeclare(
		"logs",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	common.Errlog(err, "declare exchange fail")

	q, err := ch.QueueDeclare(
		"send2mysql",
		true,
		false,
		false,
		false,
		nil,
	)
	common.Errlog(err, "declare queue error happen")

	err = ch.QueueBind(
		q.Name,
		"logsRecord",
		"logs",
		false,
		nil,
	)
	common.Errlog(err, "bind exchange to queue error")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // no ack
		false,
		false,
		false,
		nil,
	)
	common.Errlog(err, "Failed msgs")

	forever := make(chan bool)
	order := new(structure.Orders)
	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, order)
			if err != nil {
				log.Printf("解析传送过来的[]byte到结构体时, 出现了错误, %v", err)
			}
			log.Printf("Received msg: %v", order)
			err = d.Ack(false)
			if err != nil {
				log.Printf("ack error: %s", err)
			}
			// 开始将redis来的订单信息同步到数据库中
			err = orders.InsertOrders(order.OrderNum, order.UserId, order.ProductId, order.PurchaseNum, order.OrderDatetime, order.Status)
			if err != nil {
				log.Printf("%s", err)
			}
		}
	}()
	log.Printf("Listening...")
	<-forever
}
