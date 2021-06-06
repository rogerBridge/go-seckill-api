package controllers

import (
	"go-seckill/internal/logconf"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/rabbitmq/common"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "controllers"})

var ch = common.Ch

// 全局变量, 存储purchase_limits
var purchaseLimit = make(map[string]*structure.PurchaseLimits)

// User is a type to be exported
type User struct {
	userID string
}

type Jwt struct {
	Username     string    `json:"username"`
	Jwt          string    `json:"token"` // JWT token
	GenerateTime time.Time `json:"generateTime"`
	ExpireTime   time.Time `json:"expireTime"`
}
