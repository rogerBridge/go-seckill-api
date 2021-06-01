package utils

import (
	"redisplay/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "utils"})
