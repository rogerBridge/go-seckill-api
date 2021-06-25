package sender

import (
	"go-seckill/internal/rabbitmq/common"

	"github.com/streadway/amqp"
)

func Send(msg []byte, ch *amqp.Channel) error {
	err := ch.Publish(
		"order-process",
		"", // message's routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // msg persistent
			ContentType:  "text/plain",
			Body:         msg,
		},
	)
	common.Errlog(err, "publish msg to queue error")
	return err
}
