package controllers

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go_redis/auth"
	"go_redis/jsonStruct"
	"go_redis/mysql/shop/structure"
	"go_redis/mysql/shop/users"
	"go_redis/utils"
	"log"
)

// 登录, 并且获取token
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
	exist, err := users.VerifyUsers(user)
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

func Logout(ctx *fasthttp.RequestCtx) {
	utils.ResponseWithJson(ctx, 500, structure.UserLogout{Message: "系统维护"})
	return
}

// 用户注册必须提供的参数: 用户名, 密码
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
	if err!=nil {
		log.Printf("insert users error: %s\n", err)
		utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
			Code: 8401,
			Msg:  "insert to users error",
			Data: nil,
		})
		return
	}else {
		utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
			Code: 8001,
			Msg:  "register success",
			Data: nil,
		})
		return
	}
}
