package receive

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go_redis/mysql/shop/orders"
	"go_redis/mysql/shop/structure"
	"go_redis/rabbitmq/common"
	"log"
)

func Receive(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"send2mysql",
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
			err = d.Ack(false)
			if err != nil {
				log.Printf("ack error: %s", err)
			}
			log.Printf("%s", d.Body)

			err := json.Unmarshal(d.Body, order)
			if err != nil {
				log.Printf("解析传送过来的[]byte到结构体时, 出现了错误, %v", err)
			}
			log.Printf("Received msg: %v", order)
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
