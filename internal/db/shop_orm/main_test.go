package shop_orm_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var conn = &gorm.DB{}

func TestMain(m *testing.M) {
	dataSource := "root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"
	sqlDB, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gorm connect to mysql: ", dataSource)
	conn = gormDB

	// start test Main
	os.Exit(m.Run())

}
