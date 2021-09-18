package receiver

import (
	"encoding/json"
	"go-seckill/internal/db"
	"go-seckill/internal/db/shop_orm"

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
	good := new(shop_orm.Good)
	// receives order msg from queue: order
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
			// 解析order信息, 将库存相关的从mysql对应数据库中减去
			err = good.UpdateGoodByProductIDandPurchaseNum(db.Conn2, order.ProductID, order.PurchaseNum)
			if err != nil {
				logger.Warnf("更新商品库存时出现错误: %v", err)
			}
			err = order.CreateOrder(db.Conn2)
			//err = orders.InsertOrders(order.OrderNum, order.UserId, order.ProductId, order.PurchaseNum, order.OrderDatetime, order.Status)
			if err != nil {
				logger.Warnf("将mqtt接收到的消息同步到order表格时出现错误: %s", err)
			}
		}
	}()
	logger.Infof("Listening incoming mqtt info ...")
	<-forever
}
