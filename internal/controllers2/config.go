package controllers2

import (
	"go-seckill/internal/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "controllers"})
