package main

import (
	"go-seckill/build/rabbitmqReceiver/receiver"
	"go-seckill/internal/rabbitmq/common"
)

func main() {
	receiver.Receive(common.Ch)
}
