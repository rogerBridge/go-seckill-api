package auth

import (
	"redisplay/logconf"
	"time"

	"github.com/sirupsen/logrus"
)

// server side sign token need secret
// stateless token
var secret = "1hXNV1rlgoEoT9U9gWqSmyYS9G1"

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "auth"})

// 过期时间, 自定义, 我这里配置的是: 7天
var ExpireDuration = 7 * 24 * time.Hour
