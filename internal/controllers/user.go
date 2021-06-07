package controllers

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/auth"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/mysql/shop/users"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"time"

	"github.com/valyala/fasthttp"
)

// Login 登录, 生成token, 并放入tokenRedis中
func Login(ctx *fasthttp.RequestCtx) {
	//var user structure.UserLogin
	user := new(structure.UserLogin)
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		logger.Warnf("Login: Unmarshal []byte to struct error message %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "bad params",
			Data: nil,
		})
		return
	}
	exist, err := users.VerifyUsers(user)
	if err != nil {
		logger.Warnf("Login: While verify user: %s error message: %v", user.Username, err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  fmt.Sprintf("While verify user: %v error message: %v", user, err),
			Data: nil,
		})
		return
	}
	if exist == 1 {
		token, err := auth.GenerateToken(user)
		if err != nil {
			logger.Warnf("Login: While verify user: %v error message: %v", user, err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  fmt.Sprintf("While verify user: %v error message: %v", user, err),
				Data: nil,
			})
			return
		}
		logger.Infof("Login: User: %v login success", user)
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8001,
			Msg:  "login success",
			Data: Jwt{
				Username:     user.Username,
				Jwt:          token,
				GenerateTime: time.Now(),
				ExpireTime:   time.Now().Add(auth.ExpireDuration),
			},
		})
	} else {
		logger.Warnf("Login: username or password error")
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "用户名或密码错误",
			Data: nil,
		})
	}
}

// Logout 删除token
func Logout(ctx *fasthttp.RequestCtx) {
	// get username from token string
	tokenStr := string(ctx.Request.Header.Peek("Authorization"))
	//
	tokenInfo, err := auth.ParseToken(tokenStr)
	if err != nil {
		logger.Warnf("Logout: user: %s logout error message: %v", tokenInfo.Username, err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  fmt.Sprintf("user: %v logout error message: %v", tokenInfo, err),
			Data: nil,
		})
		return
	}
	username := tokenInfo.Username
	// tokenRedis
	redisconn := redisconf.Pool2.Get()
	defer redisconn.Close()

	_, err = redisconn.Do("del", "token:"+username)
	if err != nil {
		logger.Warnf("Logout: user: %s delete self token error message: %v", username, err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  fmt.Sprintf("user: %s delete self token error message: %v", username, err),
			Data: nil,
		})
		return
	}
	logger.Infof("Logout: user %s login success", username)
	utils.ResponseWithJson(ctx, 200, structure.UserLogout{Message: "logout successful"})
}

// Register 用户注册必须提供的参数: 用户名, 密码, 邮箱
func Register(ctx *fasthttp.RequestCtx) {
	user := new(structure.UserRegister)
	err := json.Unmarshal(ctx.Request.Body(), user)
	if err != nil {
		logger.Warnf("Register: user: %v register json unmarshal error message: %v", user, err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "bad params",
			Data: nil,
		})
		return
	}
	// 这里应该填写用户名, 密码, 邮箱的校验机制, 应该用正则表达式, 暂时偷懒
	if user.Email == "" || user.Username == "" || user.Password == "" {
		logger.Warnf("Register: user: %v must have username, password, email", user)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "incomplete register info",
			Data: nil,
		})
		return
	}
	// 首先查找数据库中是否存在这个用户
	_, err = users.VerifyIfUserExist(user)
	if err != nil {
		logger.Warnf("User or Email is Exist or Other Error")
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "username or email existed or Other Error",
			Data: nil,
		})
		return
	}
	// 设置默认birthday值
	if user.Birthday == "" {
		user.Birthday = "2006-01-02 13:04:05"
	}
	err = users.InsertUsers(user)
	if err != nil {
		logger.Warnf("Register: Insert user: %v error message: %v", user, err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8401,
			Msg:  fmt.Sprintf("Insert user: %v error message: %v", user, err),
			Data: nil,
		})
		return
	}
	logger.Infof("Register: User: %v register success", user)
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8001,
		Msg:  "register success",
		Data: nil,
	})
}
