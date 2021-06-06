package receiver

import (
	"encoding/json"
	"go-seckill/internal/mysql/shop/orders"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/rabbitmq/common"

	"github.com/streadway/amqp"
)

// 从特定的channel里面接收信息, 然后处理
func Receive(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"sendToMysql", // queue name
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
				logger.Warnf("ack error: %s", err)
			}
			// log.Printf("%s", d.Body)

			err := json.Unmarshal(d.Body, order)
			if err != nil {
				logger.Warnf("解析传送过来的[]byte到结构体时, 出现了错误, %v", err)
			}
			logger.Warnf("Received msg: %+v", order)
			// 开始将redis来的订单信息同步到数据库中
			err = orders.InsertOrders(order.OrderNum, order.UserId, order.ProductId, order.PurchaseNum, order.OrderDatetime, order.Status)
			if err != nil {
				logger.Warnf("%s", err)
			}
		}
	}()
	logger.Warnf("Listening incoming mqtt info ...")
	<-forever
}
