package main

import (
	"go-seckill/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "pressuretest"})
