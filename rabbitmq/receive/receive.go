package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"go_redis/rabbitmq/common"
	"log"
	"time"
)

func Receive() {
	//
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
		"",
		false,
		false,
		true,
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
	go func() {
		for d := range msgs {
			dotCount := bytes.Count(d.Body, []byte("."))
			time.Sleep(time.Duration(dotCount) * time.Second)
			log.Printf("Received msg: %s", d.Body)
			err := d.Ack(false)
			if err != nil {
				log.Printf("ack error: %s", err)
			}
		}
	}()
	log.Printf("Listening...")
	<-forever
}
