package users

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"redisplay/mysql"
	"redisplay/mysql/shop/structure"

	_ "github.com/go-sql-driver/mysql"
)

var salt = "1hXNV1rlgoEoT9U9gWqSmyYS9G1"

// 向users table中insert数值
func InsertUsers(u *structure.UserRegister) error {
	_, err := mysql.Conn.Exec("insert users (name, passwd, sex, birthday, address, email) values (?, ?, ?, ?, ?, ?)", u.Username, GetMD5sum(u.Password), u.Sex, u.Birthday, u.Address, u.Email)
	if err != nil {
		log.Printf("insert users error: %s\n", err)
		return err
	}
	return nil
}

// 更新user信息
func UpdateUsers(u *structure.UserLogin) error {
	_, err := mysql.Conn.Exec("update users set passwd=? where name=?", GetMD5sum(u.Password), u.Username)
	if err != nil {
		log.Printf("update users passwd error: %s\n", err)
		return err
	}
	return nil
}

// 验证登录用户名和密码是否存在
func VerifyUsers(u *structure.UserLogin) (int, error) {
	var isExist int
	row := mysql.Conn.QueryRow("select 1 from users where name=? and passwd=?", u.Username, GetMD5sum(u.Password))
	err := row.Scan(&isExist)
	if err != nil {
		log.Printf("verify username and password error: %v\n", err)
		return 0, err
	}
	if isExist != 1 {
		log.Printf("user not exist \n")
		return isExist, nil
	}
	//log.Printf("user: %v exist\n", u)
	return isExist, nil
}

// 将text string转为 md5 string
func GetMD5sum(text string) string {
	hash := md5.Sum([]byte(text + salt))
	return hex.EncodeToString(hash[:])
}
