package shop_orm_test

import (
	"log"
	"testing"

	"gorm.io/gorm"
)

var conn = &gorm.DB{}

func TestMain(m *testing.M) {
	dataSource := "root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println(dataSource)
}
