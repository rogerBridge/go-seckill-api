package redisconf

import (
	"go-seckill/internal/logconf"
	"go-seckill/internal/rabbitmq/common"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "redismethods"})

var ch = common.Ch
