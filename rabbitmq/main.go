package main

import (
	"redisplay/rabbitmq/common"
	"redisplay/rabbitmq/receive"
)

func main() {
	receive.Receive(common.Ch)
}
