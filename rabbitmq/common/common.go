package common

import (
	"github.com/streadway/amqp"
	"log"
)

func Errlog(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}

var Ch = GetChannel()

func GetChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://root:12345678@my_rabbit:5672/root_vhost")
	Errlog(err, "Failed to connect my_rabbit")
	//defer conn.Close()

	ch, err := conn.Channel()
	Errlog(err, "Failed to establish channel")
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	Errlog(err, "declare exchange fail")

	q, err := ch.QueueDeclare(
		"send2mysql",
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
