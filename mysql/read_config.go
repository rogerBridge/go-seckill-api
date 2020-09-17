package mysql

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
)

// 定义mysql连接为全局变量
var Conn *sql.DB

func init() {
	pwd, err := os.Getwd()
	if err!=nil {
		log.Fatalf("找不到当前工作路径, %s", err)
	}
	log.Printf("Current Work Dir: %s", pwd)
	dataSource := ReadConfig(pwd+"/mysql/mysql_config.json")

	Conn = InitMysqlConn(dataSource)
}

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

func InitMysqlConn(dataSource string) *sql.DB {
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
