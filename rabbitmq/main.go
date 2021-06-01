package main

import (
	"redisplay/rabbitmq/common"
	"redisplay/rabbitmq/receiver"
)

func main() {
	receiver.Receive(common.Ch)
}
