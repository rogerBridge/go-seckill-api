package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
	"go_redis/jsonStruct"
	"go_redis/mysql/shop/structure"
	"go_redis/redis_config"
	"go_redis/utils"
	"log"
	"time"
)

// server side sign token need secret
var secret = "1hXNV1rlgoEoT9U9gWqSmyYS9G1"

// 生成符合要求的JWT token
// 要求如下: 24h后过期
func GenerateToken(user *structure.UserLogin) (string, error) {
	expireDuration := 24 * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(expireDuration).Unix(), // 24 hours expire
	})
	tokenReturn, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("error happen while generate token\n")
	}
	// 将生成的token放入redis
	redisconn := redis_config.Pool2.Get()
	defer redisconn.Close()

	_, err = redisconn.Do("set", "token:"+user.Username, tokenReturn)
	_, err = redisconn.Do("expire", "token:"+user.Username, int64(expireDuration)/1e9) // 1e9 = 1 Second
	if err != nil {
		log.Fatalln("error while set user:token", err)
	}
	return tokenReturn, nil
}

// token middleware
func MiddleAuth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// 首先, 验证header中key: authorization的值是否符合要求?
		tokenStr := string(ctx.Request.Header.Peek("Authorization"))
		tokeninfo, err := ParseToken(tokenStr)
		if err != nil {
			utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
				Code: 8400,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		username := tokeninfo.Username
		redisconn := redis_config.Pool2.Get()
		defer redisconn.Close()

		tokenFromRedis, err := redis.String(redisconn.Do("get", "token:"+username))
		if err != nil {
			log.Printf("while get token from redis, error: %+v\n", err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "while get token from redis, error",
				Data: nil,
			})
			return
		}
		if tokenFromRedis != tokeninfo.TokenString {
			log.Printf("unvalid token")
			utils.ResponseWithJson(ctx, 400, jsonStruct.CommonResponse{
				Code: 8400,
				Msg:  "while get token from redis, error",
				Data: nil,
			})
			return
		}
		handler(ctx)
	}
}

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
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, err error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return
			}
			return []byte(secret), nil // default return
		})
		// err 中包含错误
		if err != nil {
			log.Printf("token parse error: %+v\n", err)
			return tokenInfo, err
		}
		// 如果可以顺利解析, 将解析后的值分配到 tokenInfo 结构体中
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			tokenInfo.TokenString = tokenStr
			tokenInfo.Username = claims["username"].(string)
			tokenInfo.Expiration = int64(claims["exp"].(float64))
			return tokenInfo, nil
		} else {
			return tokenInfo, errors.New("token parse to struct error")
		}
	}
}
