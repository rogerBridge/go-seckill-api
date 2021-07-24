package shop_orm

import (
	"go-seckill/internal/logconf"
	"go-seckill/internal/mysql"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 初始化数据库 seckill
func InitialMysql() {
	conn := mysql.Conn2

	conn.Exec("CREATE DATABASE IF NOT EXISTS seckill")
	log.Println("executed create database seckill command")

	err := conn.AutoMigrate(&Good{})
	if err != nil {
		log.Fatalln("While migrate goods table, error: ", err)
	}

	err = conn.AutoMigrate(&PurchaseLimit{})
	if err != nil {
		log.Fatalln("While migrate purchaseLimits table, error: ", err)
	}

	err = conn.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("While migrate users table, error: ", err)
	}

	err = conn.AutoMigrate(&Order{})
	if err != nil {
		log.Fatalln("While migrate orders table, error: ", err)
	}

}

type SelfDefine struct {
	gorm.Model
	Version string `gorm:"default:v0.0.0"`
}

var conn = mysql.Conn2

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "shop_orm"})
