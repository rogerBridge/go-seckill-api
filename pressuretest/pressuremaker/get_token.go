package pressuremaker

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// 每次跑测试前, 找后台申请一个最新的token
func GetToken() (string, error) {
	url := "http://localhost:4000/user/login"
	client := &fasthttp.Client{
		Dial: func(addr string) (conn net.Conn, err error) {
			return fasthttp.DialTimeout(addr, 60*time.Second)
		},
		ReadTimeout: 30 * time.Second,
	}
	req := new(fasthttp.Request)
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	l := new(login)
	l.Username = "roger"
	l.Password = "12345678"
	reqBytes, err := json.Marshal(l)
	if err != nil {
		logger.Warnln(reqBytes, err)
	}
	req.SetBody(reqBytes)

	resp := new(fasthttp.Response)
	err = client.Do(req, resp)
	if err != nil {
		logger.Fatalln(err)
		return "", err
	}
	type Data struct {
		Username     string    `json:"username"`
		Token        string    `json:"token"`
		GenerateTime time.Time `json:"generateTime"`
		ExpireTime   time.Time `json:"expireTime"`
	}
	type LoginInfo struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data Data   `json:"data"`
	}
	loginInfo := new(LoginInfo)
	if err != nil {
		logger.Warnln(err)
		return "", err
	}
	err = json.Unmarshal(resp.Body(), loginInfo)
	if err != nil {
		logger.Fatalln(err)
	}
	token := loginInfo.Data.Token
	logger.Infoln("get token this time: ", token)
	return token, nil
}
