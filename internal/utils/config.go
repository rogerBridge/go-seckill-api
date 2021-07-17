package utils

import (
	"go-seckill/internal/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "utils"})

const API_VERSION string = "/api/v0"
