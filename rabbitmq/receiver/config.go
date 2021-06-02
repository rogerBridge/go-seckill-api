package receiver

import (
	"redisplay/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "rabbitmq-receiver"})
