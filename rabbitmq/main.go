package main

import (
	"go_redis/rabbitmq/common"
	"go_redis/rabbitmq/receive"
)

func main() {
	receive.Receive(common.Ch)
}
