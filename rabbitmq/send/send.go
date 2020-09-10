package send

import (
	"github.com/streadway/amqp"
	"go_redis/rabbitmq/common"
	"strings"
)

func Send(msg []byte) {
	conn, err := amqp.Dial("amqp://root:12345678@my_rabbit:5672/root_vhost")
	common.Errlog(err, "Failed to connect my_rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	common.Errlog(err, "Failed to establish channel")
	defer ch.Close()

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
		"logsRecord", // receive only key match msg publish
		"logs",
		false,
		nil,
	)
	common.Errlog(err, "bind exchange to queue error")

	err = ch.Publish(
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

func args(args []string) string {
	if len(args) < 2 || args[1] == "" {
		return "hello"
	} else {
		return strings.Join(args[1:], " ")
	}
}
