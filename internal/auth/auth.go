/*
Package auth is a package provide go-seckill API auth service
*/
package auth

import (
	"errors"
	"fmt"
	"go-seckill/internal/easyjsonprocess"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
)

// 生成符合要求的JWT token
func GenerateToken(user *structure.UserLogin) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ExpireDuration).Unix(),
		Id:        user.Username,
	})
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
func MiddleAuth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// 首先, 验证header中key: authorization的值是否符合要求?
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
		// 这里存在小问题
		// 显示: 帐号已经登录, 请在别处退出后, 再重新登录
		// 应该是: 把之前的token替换掉, 并提示: 帐号在别处登录, 已经被踢出, 以新的登录为主
		// 或者是可以存储多个token
		if tokenFromRedis != tokeninfo.TokenString {
			logger.Warnf("MiddleAuth: user: %v request token not equal to tokenRedis's token", tokeninfo.Username)
			utils.ResponseWithJson(ctx, 400, easyjsonprocess.CommonResponse{
				Code: 8400,
				Msg:  "提交的token与服务器缓存的token不符",
				Data: nil,
			})
			return
		}
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
}

// Parsing token:string
func ParseToken(tokenStr string) (*TokenInfo, error) {
	tokenInfo := new(TokenInfo)
	if tokenStr == "" {
		return tokenInfo, errors.New("nil tokenStr")
	} else {
		// 验证token是否可以被解析
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return
			}
			return []byte(secret), nil
		})
		// err 中包含错误
		if err != nil {
			logger.Warnf("ParseToken: Token parse error: %v", err)
			return tokenInfo, err
		}
		// 如果可以顺利解析, 将解析后的值分配到 tokenInfo 结构体中
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			tokenInfo.TokenString = tokenStr
			tokenInfo.Username = claims.Id
			tokenInfo.Expiration = claims.ExpiresAt
			return tokenInfo, nil
		} else {
			logger.Warnf("ParseToken: error happen when parse token to struct tokenInfo, error message: %v", err)
			return tokenInfo, err
		}
	}
}
