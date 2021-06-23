package controllers2

import (
	"encoding/json"
	"fmt"
	"go-seckill/internal/auth"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/mysql/shop_orm"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"time"

	"github.com/valyala/fasthttp"
)

// 用户注册
func UserRegister(ctx *fasthttp.RequestCtx) {
	user := new(shop_orm.User)
	err := json.Unmarshal(ctx.Request.Body(), user)
	if err != nil {
		logger.Warnf("Register: user: %+v register json unmarshal error message: %v", user, err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "bad params",
			Data: nil,
		})
		return
	}
	logger.Infof("unmarshal []byte to struct successful")

	tx := mysql.Conn2.Begin()
	err = user.CreateUser(tx)
	if err != nil {
		logger.Warnf("Register transaction error: %s", err.Error())
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "注册失败: " + err.Error(),
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("Register: CreateUser: %v error: %v", user, err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8200,
			Msg:  fmt.Sprintf("CreateUser: %v error: %v", user, err),
			Data: nil,
		})
		return
	}
	logger.Infof("CreateUser: %+v register success", user)
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8200,
		Msg:  "register success",
		Data: nil,
	})
}

func UserLogin(ctx *fasthttp.RequestCtx) {
	//var user structure.UserLogin
	user := new(shop_orm.User)
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		logger.Warnf("Login: Unmarshal []byte to struct error message %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "bad params",
			Data: nil,
		})
		return
	}
	logger.Infof("Login: unmarshal []byte to struct error: %v", err)

	type JWT struct {
		Username     string    `json:"username"`
		Token        string    `json:"token"`
		GenerateTime time.Time `json:"generateTime"`
		ExpireTime   time.Time `json:"expireTime"`
	}

	// userInMysql 是真实存在于mysql中的User Object
	userInMysql, exist := user.ProofCredential()
	if exist {
		token, err := auth.GenerateToken(&userInMysql)
		if err != nil {
			logger.Warnf("Login: While verify user: %v error message: %v", user, err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  fmt.Sprintf("While verify user: %v error message: %v", user, err),
				Data: nil,
			})
			return
		}
		logger.Infof("Login: User: %+v login success", user)
		utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
			Code: 8001,
			Msg:  "login success",
			Data: JWT{
				Username:     user.Username,
				Token:        token,
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

func UserLogout(ctx *fasthttp.RequestCtx) {
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
	logger.Infof("Logout: user %s logout success", username)
	utils.ResponseWithJson(ctx, 200, structure.UserLogout{Message: "logout successful"})
}

func UserUpdatePassword(ctx *fasthttp.RequestCtx) {
	// usernameBytes := ctx.Request.Header.Peek("username")
	// var sb strings.Builder
	// sb.Grow(64)
	// for _, r := range usernameBytes {
	// 	sb.WriteByte(r)
	// }
	// username := sb.String()
	username := string(ctx.Request.Header.Peek("username"))

	p := new(shop_orm.User)
	err := json.Unmarshal(ctx.Request.Body(), p)
	if err != nil {
		logger.Warnf("unmarshal password error: %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "params unmarshal error",
			Data: nil,
		})
		return
	}
	logger.Infof("unmarshal password successful")
	p.Username = username

	tx := mysql.Conn2.Begin()
	err = p.UpdateUserPassword(tx)
	if err != nil {
		logger.Warnf("UpdateUserPassword transaction error: %v", err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "UpdateUserPassword error",
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("UpdateUserPassword transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "UpdateUserPassword error",
			Data: nil,
		})
		return
	}
	logger.Infof("%+v UpdateUserPassword successful", p)
	utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
		Code: 8200,
		Msg:  "UpdateUserPassword successful",
		Data: nil,
	})
}

func UserUpdateInfo(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	p := new(shop_orm.User)
	err := json.Unmarshal(ctx.Request.Body(), p)
	if err != nil {
		logger.Warnf("unmarshal userInfo error: %v", err)
		utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
			Code: 8400,
			Msg:  "params unmarshal error",
			Data: nil,
		})
		return
	}
	logger.Infof("unmarshal userInfo successful")
	// check userinfo
	// if !p.IfUserExist() {
	// 	logger.Warnf(err.Error())
	// 	utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
	// 		Code: 8400,
	// 		Msg:  "用户不存在",
	// 		Data: nil,
	// 	})
	// 	return
	// }
	p.Username = username

	// if err := p.CheckEmail(); err != nil {
	// 	logger.Warnf(err.Error())
	// 	utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
	// 		Code: 8400,
	// 		Msg:  err.Error(),
	// 		Data: nil,
	// 	})
	// }

	tx := mysql.Conn2.Begin()
	err = p.UpdateUserInfo(tx)
	if err != nil {
		logger.Warnf("UpdateUserInfo transaction error: %v", err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Warnf("UpdateUserInfo transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
			Code: 8500,
			Msg:  "UpdateUserInfo commit error",
			Data: nil,
		})
		return
	}
	logger.Infof("UpdateUserInfo transaction commit successful")
	utils.ResponseWithJson(ctx, 200, easyjsonprocess.CommonResponse{
		Code: 8200,
		Msg:  "UpdateUserInfo successful",
		Data: nil,
	})
}
