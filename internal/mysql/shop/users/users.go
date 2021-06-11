package users

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop/structure"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 向users table中insert数值
func InsertUsers(u *structure.UserRegister) error {
	_, err := mysql.Conn.Exec("insert users (name, passwd, sex, birthday, address, email) values (?, ?, ?, ?, ?, ?)", u.Username, GetSha256sum(u.Password), u.Sex, u.Birthday, u.Address, u.Email)
	if err != nil {
		log.Printf("insert users error: %s\n", err)
		return err
	}
	return nil
}

// 更新user信息
func UpdateUsers(u *structure.UserLogin) error {
	_, err := mysql.Conn.Exec("update users set passwd=? where name=?", GetSha256sum(u.Password), u.Username)
	if err != nil {
		log.Printf("update users passwd error: %s\n", err)
		return err
	}
	return nil
}

// 验证登录用户名和密码是否存在
func FindUserIfExist(u *structure.UserLogin) (int, error) {
	var isExist int
	row := mysql.Conn.QueryRow("select 1 from users where name=? and passwd=?", u.Username, GetSha256sum(u.Password))
	err := row.Scan(&isExist)
	if err != nil {
		log.Printf("row scan error happend : %v\n", err)
		return 0, err
	}
	if isExist != 1 {
		log.Printf("user not exist \n")
		return isExist, nil
	}
	//log.Printf("user: %v exist\n", u)
	return isExist, nil
}

// 验证用户注册时, 用户名, 邮箱是否重复
func VerifyIfUserExist(u *structure.UserRegister) (int, error) {
	var count int
	row := mysql.Conn.QueryRow("select count(*) from users")
	err := row.Scan(&count)
	if err != nil {
		log.Printf("row scan error: %v\n", err)
		return 0, err
	}
	if count == 0 {
		return 0, nil
	}
	var ifUserExist int
	row = mysql.Conn.QueryRow("select count(*) from users where name=? or email=?", u.Username, u.Email)
	err = row.Scan(&ifUserExist)
	if err != nil {
		log.Printf("row scan error :%v\n", err)
		return 0, err
	}
	if ifUserExist == 1 {
		log.Printf("user exist! \n")
		return 1, fmt.Errorf("user exist")
	}
	return ifUserExist, nil
}

// 将text string转为 base64(sha256(passwd+salt))
func GetSha256sum(text string) string {
	v := sha256.Sum256([]byte(text + salt))
	v1 := v[:]
	return base64.URLEncoding.EncodeToString(v1)
}
