package receiver

import (
	"encoding/json"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop_orm"

	"go-seckill/internal/rabbitmq/common"

	"github.com/streadway/amqp"
)

// receive msg from channel, then process
func Receive(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"order", // queue name
		"",
		false, // no ack
		false,
		false,
		false,
		nil,
	)
	common.Errlog(err, "Failed msgs")
	forever := make(chan bool)
	order := new(shop_orm.Order)
	go func() {
		for d := range msgs {
			err := d.Ack(false)
			if err != nil {
				logger.Warnf("ack error: %s", err)
			}
			// log.Printf("%s", d.Body)

			err = json.Unmarshal(d.Body, order)
			if err != nil {
				logger.Warnf("解析传送过来的[]byte到结构体时, 出现了错误, %v", err)
			}
			logger.Infof("Received msg: %+v", order)
			// 开始将redis来的订单信息同步到数据库中
			// err = order.UpdateOrder(mysql.Conn2)
			err = order.CreateOrder(mysql.Conn2)
			//err = orders.InsertOrders(order.OrderNum, order.UserId, order.ProductId, order.PurchaseNum, order.OrderDatetime, order.Status)
			if err != nil {
				logger.Warnf("将mqtt接收到的消息同步到order表格时出现错误: %s", err)
			}
		}
	}()
	logger.Infof("Listening incoming mqtt info ...")
	<-forever
}
