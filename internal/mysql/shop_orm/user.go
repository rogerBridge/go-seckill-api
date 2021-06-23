package shop_orm

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type User struct {
	SelfDefine
	Username string    `gorm:"username" json:"username"`
	Password string    `gorm:"password" json:"password"`
	Group    string    `gorm:"default:user" json:"group"`
	Sex      string    `gorm:"sex" json:"sex"`
	Birthday time.Time `gorm:"birthday" json:"birthday"`
	Address  string    `gorm:"address" json:"address"`
	Email    string    `gorm:"email" json:"email"`
}

type UserJson struct {
	SelfDefine
	Username string `json:"username"`
	Password string `json:"password"`
	Group    string `json:"group"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// base64(sha256sum(password+salt))
func passwordEncrypt(password string) string {
	salt := "a952a114-8a87-4617-8285-19f000e41c9e"
	s := (sha256.Sum256([]byte(password + salt)))
	ss := s[:]
	return base64.StdEncoding.EncodeToString(ss)
}

func (u *User) CreateUser(tx *gorm.DB) error {
	// 检测用户注册信息是否符合规范
	if err := u.CheckEmailFormat(); err != nil {
		return err
	}
	if err := u.CheckEmailIsUnique(); err != nil {
		return err
	}
	// 检测系统中是否存在此用户
	if u.IfUserExist() {
		return fmt.Errorf("用户已存在, 无法新建")
	}
	// 需要将密码切换为sha256sum+salt的形式
	u.Password = passwordEncrypt(u.Password)
	if err := tx.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) QueryUsers() ([]*User, error) {
	users := make([]*User, 128)
	if err := conn.Model(&User{}).Find(users).Error; err != nil {
		return users, err
	}
	return users, nil
}

// 更新用户信息(除密码之外)
func (u *User) UpdateUserInfo(tx *gorm.DB) error {
	if err := u.CheckEmailFormat(); err != nil {
		return err
	}
	if err := tx.Model(&User{}).Where("username=?", u.Username).Updates(User{Email: u.Email, Sex: u.Sex, Birthday: u.Birthday, Address: u.Address}).Error; err != nil {
		if err != nil {
			log.Println("UpdateUserInfo error: ", err)
			return err
		}
	}
	return nil
}

// 更新用户密码
func (u *User) UpdateUserPassword(tx *gorm.DB) error {
	if !u.CheckPasswordValid() {
		return fmt.Errorf("密码不符合要求")
	}
	if err := tx.Model(&User{}).Where("username=?", u.Username).Update("password", passwordEncrypt(u.Password)).Error; err != nil {
		log.Println("UpdateUserPassword error: ", err)
		return err
	}
	return nil
}

// 传入的指定User是否存在于users table, 检测username是否存在, 注册user的时候使用
// username和email都必须唯一
func (u *User) IfUserExist() bool {
	var result User
	if err := conn.Model(&User{}).Where("username=? OR email=?", u.Username, u.Email).First(&result).Error; err == nil {
		if result.Username != "" || result.Email != "" {
			return true
		}
	}
	return false
}

// 检查用户名和密码是否符合
// 如果符合, 返回对应行的数据
func (u *User) ProofCredential() (User, bool) {
	var result User
	if err := conn.Model(&User{}).Where("username=? AND password=?", u.Username, passwordEncrypt(u.Password)).First(&result).Error; err == nil {
		if result.Username != "" && result.Password != "" {
			return result, true
		}
	}
	return result, false
}

// 检查用户信息是否符合要求
func (u *User) CheckEmailFormat() error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(u.Email) < 3 || len(u.Email) > 255 {
		return fmt.Errorf("邮箱地址长度不符合要求")
	}
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("邮箱地址格式不符合要求")
	}
	return nil
}

func (u *User) CheckEmailIsUnique() error {
	checkEmailIsUnique := new(User)
	if err := conn.Model(&User{}).Where("email=?", u.Email).First(checkEmailIsUnique).Error; err == nil {
		if checkEmailIsUnique.Email != "" {
			return fmt.Errorf("邮箱已存在")
		}
	}
	return nil
}

func (u *User) CheckPasswordValid() bool {
	return len(u.Password) >= 8
}
