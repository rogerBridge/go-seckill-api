package mysql

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go-seckill/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Conn 定义mysql连接为全局变量
var Conn = InitMysqlConn()

// Conn2 定义gorm连接mysql的全局变量
var Conn2 = InitMysqlConn2()

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

type MysqlConn struct {
	Username string
	Password string
	IP       string
	Port     string
	Database string
	Other    string
}

func DataSourceMysqlConn(databaseName string) string {
	//pwd, err := os.Getwd()
	//if err != nil {
	//	log.Fatalf("找不到当前工作路径, %s", err)
	//}
	//log.Printf("当前工作路径是: %s\n", pwd)

	//viper.SetConfigFile("./config/mysql_config.json")
	viper.SetConfigType("json")
	err := viper.ReadConfig(bytes.NewBuffer(config.ReadConfig("mysql_config.json")))
	if err != nil {
		log.Fatalf("当使用viper读取mysql配置的时候出错, 程序崩溃, 错误信息: %v", err)
	}
	dbInstance := viper.GetStringMapString(databaseName)

	mysqlconn := new(MysqlConn)
	mysqlconn.Username = dbInstance["username"]
	mysqlconn.Password = dbInstance["password"]
	mysqlconn.IP = dbInstance["ip"]
	mysqlconn.Port = dbInstance["port"]
	mysqlconn.Database = dbInstance["database"]
	mysqlconn.Other = "?charset=utf8mb4&parseTime=True&loc=Local"

	dataSource := mysqlconn.Username + ":" + mysqlconn.Password + "@tcp(" + mysqlconn.IP + ":" + mysqlconn.Port + ")/" + mysqlconn.Database + mysqlconn.Other
	log.Println("从viper读取到的mysql的配置是:", dataSource)
	return dataSource
}

// InitMysqlConn ...
// 使用viper初始化mysql实例
func InitMysqlConn() *sql.DB {
	dataSource := DataSourceMysqlConn("db")
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("conn establish error: %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("conn establish error: %v\n", err)
	}
	db.SetMaxIdleConns(10)
	log.Println("mysql go-libray connect to mysql: ", dataSource)
	// db.SetMaxOpenConns(100)
	// db.SetConnMaxLifetime(time.Hour)
	return db
}

// 这里是通过gorm的方式连接数据库
func InitMysqlConn2() *gorm.DB {
	dataSource := DataSourceMysqlConn("db2")
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gorm connect to mysql: ", dataSource)
	return db
}
