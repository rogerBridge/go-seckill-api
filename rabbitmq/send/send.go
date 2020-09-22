package send

import (
	"github.com/streadway/amqp"
	"go_redis/rabbitmq/common"
)

func Send(msg []byte, ch *amqp.Channel) {
	err := ch.Publish(
		"logs",
		"logsRecord", // message's routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // msg persistent
			ContentType:  "text/plain",
			Body:         msg,
		},
	)
	common.Errlog(err, "publish msg to queue error")
}
