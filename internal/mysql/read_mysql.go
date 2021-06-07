package mysql

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// Conn 定义mysql连接为全局变量
var Conn *sql.DB = InitMysqlConn()

// 可以不使用init, 就坚决不使用init

// // ReadConfig 读取mysql数据库的设置, 并输出: "username:password@tcp(IPaddress:port)/database+params" format
// func ReadConfig(fileName string) string {
// 	type mysqlConfig struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 		IP       string `json:"ip"`
// 		Port     string `json:"port"`
// 		Database string `json:"database"`
// 	}
// 	f, err := os.Open(fileName)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	var m mysqlConfig
// 	err = json.NewDecoder(f).Decode(&m)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	params := "?parseTime=true"
// 	return m.Username + ":" + m.Password + "@tcp(" + m.IP + ":" + m.Port + ")/" + m.Database + params
// }

// InitMysqlConn ...
// 使用viper初始化mysql实例
func InitMysqlConn() *sql.DB {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("找不到当前工作路径, %s", err)
	}
	log.Printf("当前工作路径是: %s\n", pwd)

	viper.SetConfigFile("./config/mysql_config.json")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("当使用viper读取mysql配置的时候出错, 程序崩溃, 错误信息: %v", err)
	}
	dbInstance := viper.GetStringMapString("db")
	username := dbInstance["username"]
	password := dbInstance["password"]
	ip := dbInstance["ip"]
	port := dbInstance["port"]
	database := dbInstance["database"]
	dataSource := username + ":" + password + "@tcp(" + ip + ":" + port + ")/" + database + "?parseTime=true"
	log.Println("从viper读取到的mysql的配置是:", dataSource)
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
