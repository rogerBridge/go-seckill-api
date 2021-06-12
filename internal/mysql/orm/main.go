package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SelfDefine struct {
	gorm.Model
	Version string `gorm:"default:v0.0.0"`
}

type Product struct {
	SelfDefine // 注意, 这里之后需要使用阿里标准table: primarykey, version, is_delete, gmtCreate, gmtUpdate,
	Code       string
	Price      uint // 单位: 分
}

type User struct {
	SelfDefine
	Name     string
	Age      int
	Birthday time.Time
	TestItem string
}

var db, errConnectToSqlite3 = gorm.Open(sqlite.Open("/home/leo/Source/sqlite3/test.db"), &gorm.Config{})

func start() {
	// connect to database

	// drop table is exists
	// db.Exec("DROP TABLE IF EXISTS users")
	if errConnectToSqlite3 != nil {
		log.Fatalln("failed to connect sqlite3")
	}

	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("While db.AutoMigrate, ", err)
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalln("While db.AutoMigrate, ", err)
	}
}

func main() {
	start()
	type u struct {
		name string
	}
	fmt.Printf("%x\n", []byte("安瑞峰"))
	// create()
	// query()
}
