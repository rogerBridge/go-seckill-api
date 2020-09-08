package mysql

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
)

// 读取mysql数据库的设置
func ReadConfig(fileName string) string {
	type mysqlConfig struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Ip       string `json:"ip"`
		Port     string `json:"port"`
		Database string `json:"database"`
	}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	var m mysqlConfig
	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		log.Fatalln(err)
	}
	params := "?parseTime=true"
	return m.Username + ":" + m.Password + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.Database + params
}

var dataSource string = ReadConfig("mysql/mysql_config.json")

func InitMysqlConn() *sql.DB {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("conn establish error: %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("conn establish error: %v\n", err)
	}
	return db
}

// 定义一个全局变量, 方便复用
var Conn = InitMysqlConn()