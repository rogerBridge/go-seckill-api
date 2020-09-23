package controllers

import (
	"encoding/json"
	"go_redis/auth"
	"go_redis/jsonStruct"
	"go_redis/mysql/shop/structure"
	"go_redis/mysql/shop/users"
	"go_redis/redis_config"
	"go_redis/utils"
	"log"

	"github.com/valyala/fasthttp"
)

// Login 登录, 并且获取token
func Login(ctx *fasthttp.RequestCtx) {
	//var user structure.UserLogin
	user := new(structure.UserLogin)
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "bad params",
			Data: nil,
		})
		return
	}
	exist, _ := users.VerifyUsers(user)
	//log.Printf("verify user error: %v\n", err)
	if exist == 1 {
		token, err := auth.GenerateToken(user)
		if err != nil {
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "while generate token error",
				Data: nil,
			})
			return
		}
		utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
			Code: 8001,
			Msg:  "login success",
			Data: structure.Jwt{
				Username: user.Username,
				Jwt:      token,
			},
		})
	} else {
		utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
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
		log.Printf(err.Error())
		utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
			Code: 8400,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	username := tokenInfo.Username
	redisconn := redis_config.Pool2.Get()
	defer redisconn.Close()

	_, err = redisconn.Do("del", "token:"+username)
	if err != nil {
		log.Printf("%+v\n", err.Error())
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "del token in redis error",
			Data: nil,
		})
		return
	}
	utils.ResponseWithJson(ctx, 200, structure.UserLogout{Message: "logout successful"})
}

// Register 用户注册必须提供的参数: 用户名, 密码
func Register(ctx *fasthttp.RequestCtx) {
	user := new(structure.UserRegister)
	err := json.Unmarshal(ctx.Request.Body(), user)
	if err != nil {
		log.Printf("parse request body error: %s\n", err)
		utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
			Code: 8400,
			Msg:  "bad params",
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
		log.Printf("insert users error: %s\n", err)
		utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
			Code: 8401,
			Msg:  "insert to users error",
			Data: nil,
		})
		return
	}
	utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "register success",
		Data: nil,
	})
}
