/*
Package auth is a package provide go-seckill API auth service
*/
package auth

import (
	"errors"
	"fmt"
	"go-seckill/internal/db/shop_orm"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
)

type MyCustomClaims struct {
	Group    string `json:"group"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成符合要求的JWT token
func GenerateToken(user *shop_orm.User) (string, error) {
	// 自定义的token
	claims := MyCustomClaims{
		user.Group,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenReturn, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Fatalf("GenerateToken: crash when generate token: %v\n", err)
	}
	// 将生成的token放入tokenRedis
	redisconn := redisconf.Pool2.Get()
	defer redisconn.Close()

	_, err = redisconn.Do("set", "token:"+user.Username, tokenReturn)
	if err != nil {
		logger.Fatalf("GenerateToken: crash when set user: %v's token, error msg: %v", user.Username, err)
	}
	_, err = redisconn.Do("expire", "token:"+user.Username, int64(ExpireDuration)/1e9) // 1e9 = 1 Second
	if err != nil {
		logger.Fatalf("GenerateToken: crash when expire user: %v's token, error msg: %v", user.Username, err)
	}
	return tokenReturn, nil
}

// MiddleAuth 是一个request前置处理器, 可以验证请求的合法性
// 只有token值合格的时候才放行请求到下一个处理器
// 这里可以做: 根据URI和group之间的关系, 做权限鉴定
func MiddleAuth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// 首先, 验证header中key: authorization的值是否符合要求?
		// 这里可以根据path判断用户是否有访问这个API的权利
		thisURI := string(ctx.URI().Path())
		tokenStr := string(ctx.Request.Header.Peek("Authorization"))
		tokeninfo, err := ParseToken(tokenStr)
		if err != nil {
			utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
				Code: 8400,
				Msg:  fmt.Sprintf("While parse token, error happen, error message is: %v", err),
				Data: nil,
			})
			return
		}
		group := tokeninfo.Group
		if !URIauthorityManage(group, thisURI) {
			logger.Warnf("MiddleAuth: 您没有访问此uri的权限")
			utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
				Code: 8401,
				Msg:  "您没有访问此uri的权限",
				Data: nil,
			})
			return
		}

		username := tokeninfo.Username
		// tokenRedis
		redisconn := redisconf.Pool2.Get()
		defer redisconn.Close()

		tokenFromRedis, err := redis.String(redisconn.Do("get", "token:"+username))
		if err != nil {
			logger.Warnf("MiddleAuth: While user: %v getting token from tokenRedis, error message: %v", tokeninfo.Username, err)
			utils.ResponseWithJson(ctx, 500, easyjsonprocess.CommonResponse{
				Code: 8500,
				Msg:  "While getting token from tokenRedis, error",
				Data: nil,
			})
			return
		}
		if tokenFromRedis != tokeninfo.TokenString {
			logger.Warnf("MiddleAuth: user: %v request token not equal to tokenRedis's token", tokeninfo.Username)
			utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
				Code: 8400,
				Msg:  "提交的token与服务器缓存的token不符",
				Data: nil,
			})
			return
		}
		// 将从token处解析到的username加到request header上, 然后发往目标api, 目标API可以直接拿到header上面的键值对
		ctx.Request.Header.Set("username", username)
		// 通过了上面的考验, 请求终于来到了handler的手上
		handler(ctx)
	}
}

/*
	TokenInfo 是解析tokenString之后的一系列信息, 例如: 用户名, 过期时间etc, 放入一个结构体中
*/
type TokenInfo struct {
	TokenString string `json:"tokenString"`
	Username    string `json:"username"`
	Expiration  int64  `json:"expiration"`
	Group       string `json:"group"`
}

// Parsing token:string
func ParseToken(tokenStr string) (*TokenInfo, error) {
	tokenInfo := new(TokenInfo)
	if tokenStr == "" {
		return tokenInfo, errors.New("nil tokenStr")
	} else {
		// 验证token是否可以被解析
		token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(secret), nil
		})
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			// logger.Infof("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
			tokenInfo.TokenString = tokenStr
			tokenInfo.Username = claims.Username
			tokenInfo.Group = claims.Group
			tokenInfo.Expiration = claims.ExpiresAt
			return tokenInfo, err
		} else {
			logger.Warnf("error: %v", err)
			return tokenInfo, err
		}
	}
}

// URI权限管理
func URIauthorityManage(group string, uri string) bool {
	if group == "user" {
		if _, ok := userURI[uri]; ok {
			return true
		}
		return false
	}
	if group == "admin" {
		if _, ok := adminURI[uri]; ok {
			return true
		}
		return false
	}
	return false
}

const ApiVersion = utils.API_VERSION

var userURI = map[string]struct{}{
	ApiVersion + "/user/logout":          {},
	ApiVersion + "/user/updatePassword":  {},
	ApiVersion + "/user/updateInfo":      {},
	ApiVersion + "/user/order/buy":       {},
	ApiVersion + "/user/order/cancelBuy": {},
	ApiVersion + "/goodList":             {},
	ApiVersion + "/user/orders":          {},
}

var adminURI = map[string]struct{}{
	ApiVersion + "/admin/createPurchaseLimit": {},
	ApiVersion + "/admin/queryPurchaseLimits": {},
	ApiVersion + "/admin/updatePurchaseLimit": {},
	ApiVersion + "/admin/deletePurchaseLimit": {},

	ApiVersion + "/admin/goodCreate": {},
	ApiVersion + "/admin/goodUpdate": {},
	ApiVersion + "/admin/goodDelete": {},
	ApiVersion + "/goodList":         {},
}
