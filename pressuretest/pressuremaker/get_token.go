package pressuremaker

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// 每次跑测试前, 找后台申请一个最新的token
func GetToken() (string, error) {
	url := "http://localhost:4000/login"
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
	l.Username = "fenrmen"
	l.Password = "12345678"
	reqBytes, _ := json.Marshal(l)
	req.SetBody(reqBytes)

	resp := new(fasthttp.Response)
	err := client.Do(req, resp)
	if err != nil {
		log.Println(err)
		return "", err
	}
	type Data struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	type LoginInfo struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data Data   `json:"data"`
	}
	loginInfo := new(LoginInfo)
	err = json.Unmarshal(resp.Body(), loginInfo)
	if err != nil {
		log.Println(err)
		return "", err
	}
	token := loginInfo.Data.Token
	log.Println("get token this time: ", token)
	return token, nil
}
