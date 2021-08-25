package pressuremaker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type data struct {
	Username     string    `json:"username"`
	Token        string    `json:"token"`
	GenerateTime time.Time `json:"generateTime"`
	ExpireTime   time.Time `json:"expireTime"`
}
type loginInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data data   `json:"data"`
}

type tokenInfo struct {
	Username string
	Token    string
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 每次跑测试前, 找后台申请一个最新的token
func (u *UserLogin) GetLoginInfo(w *sync.WaitGroup, tokenChan chan string) (*loginInfo, error) {
	url := "http://127.0.0.1:4000/api/v0/user/login"
	client := FastHttpClient
	req := new(fasthttp.Request)
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)
	logger.Println(u.Username, u.Password)

	reqBytes, err := json.Marshal(u)
	if err != nil {
		logger.Fatalln(reqBytes, err)
	}
	req.SetBody(reqBytes)

	resp := new(fasthttp.Response)
	loginInfo := new(loginInfo)
	err = client.Do(req, resp)
	if err != nil {
		logger.Warnln(err)
		//logger.Fatalln(err)
		return loginInfo, err
	}
	err = json.Unmarshal(resp.Body(), loginInfo)
	if err != nil {
		logger.Warnln(err)
		//logger.Fatalln(err)
		return loginInfo, err
	}
	logger.Debugln(string(resp.Body()))
	tokenChan <- loginInfo.Data.Token
	w.Done()
	return loginInfo, nil
}

// 每次跑测试前, 找后台申请一个最新的token, 单个goroutine
func (u *UserLogin) GetLoginInfoSingle() (*loginInfo, error) {
	url := "http://127.0.0.1:4000/api/v0/user/login"
	client := FastHttpClient
	req := new(fasthttp.Request)
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)

	reqBytes, err := json.Marshal(u)
	if err != nil {
		logger.Fatalln(reqBytes, err)
	}
	req.SetBody(reqBytes)

	resp := new(fasthttp.Response)
	err = client.Do(req, resp)
	if err != nil {
		logger.Fatalln(err)
		return &loginInfo{}, err
	}
	loginInfo := new(loginInfo)
	err = json.Unmarshal(resp.Body(), loginInfo)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Debugln(string(resp.Body()))
	return loginInfo, nil
}

// single goroutine 获取用户token
func GetTokenListSingle() ([]tokenInfo, error) {
	users := make([]*UserLogin, ConcurrentNum)
	for i := 0; i < ConcurrentNum; i++ {
		users[i] = &UserLogin{
			Username: "test" + strconv.Itoa(i),
			Password: "12345678",
		}
	}
	t0 := time.Now()
	var loginInfo *loginInfo
	var err error
	//tokenList := make([]string, 0, 10000)
	tokenList := make([]tokenInfo, 0, ConcurrentNum)
	for i := 0; i < ConcurrentNum; i++ {
		loginInfo, err = users[i].GetLoginInfoSingle()
		if err != nil {
			logger.Printf("While get user token, error: %s", err)
			return tokenList, err
		}
		if loginInfo.Data.Token != "" {
			tokenList = append(tokenList, tokenInfo{
				Username: loginInfo.Data.Username,
				Token:    loginInfo.Data.Token,
			})
		}
	}
	t1 := time.Since(t0)
	fmt.Printf("获取token用时: %dms, 获取token总个数: %d\n", t1.Milliseconds(), len(tokenList))
	return tokenList, nil
}

// 得到 []*loginInfo
func GetTokenList() {
	users := make([]*UserLogin, ConcurrentNum)
	for i := 0; i < ConcurrentNum; i++ {
		users[i] = &UserLogin{
			Username: "test" + strconv.Itoa(i),
			Password: "12345678",
		}
	}

	var w sync.WaitGroup
	tokenChan := make(chan string, ConcurrentNum)

	for i := 0; i < ConcurrentNum; i++ {
		w.Add(1)
		go users[i].GetLoginInfo(&w, tokenChan)
	}

	w.Wait()
	close(tokenChan)

	countNull := 0
	tokenList := make([]string, 0, ConcurrentNum)
	for v := range tokenChan {
		if v == "" {
			countNull += 1
		}
		tokenList = append(tokenList, v)
	}
	fmt.Println(countNull)
}

// Register single user register
func (u *User) Register() error {
	url := "http://127.0.0.1:4000/api/v0/user/register"
	client := FastHttpClient

	req := new(fasthttp.Request)
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)

	reqBytes, err := json.Marshal(u)
	if err != nil {
		logger.Warnln(reqBytes, err)
	}
	req.SetBody(reqBytes)

	resp := new(fasthttp.Response)
	err = client.Do(req, resp)
	if err != nil {
		logger.Fatalln(err)
		//return err
	}
	logger.Infoln(string(resp.Body()))
	return nil
}

// RegisterUsers batch users register
func RegisterUsers() error {
	users := make([]*User, ConcurrentNum)
	for i := 0; i < ConcurrentNum; i++ {
		users[i] = &User{
			Username: "test" + strconv.Itoa(i),
			Password: "12345678",
			Email:    "test" + strconv.Itoa(i) + "@gmail.com",
			Birthday: "2006-01-02T15:04:05+08:00",
		}
		err := users[i].Register()
		if err != nil {
			logger.Fatalln(err)
			return err
		}
	}
	return nil
}
