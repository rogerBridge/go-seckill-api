/*
Test User http api
*/

package controllers2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ReqUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespUserLogin struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type TestItemType struct {
	userLogin         ReqUserLogin
	userLoginResponse RespUserLogin
	hasError          bool
}

func TestUserLogin(t *testing.T) {
	tests := []TestItemType{
		{
			userLogin: ReqUserLogin{
				Username: "roger1",
				Password: "12345678",
			},
			userLoginResponse: RespUserLogin{
				Code: 8001,
				Msg:  "login success",
				Data: nil,
			},
			hasError: false,
		},
		{
			userLogin: ReqUserLogin{
				Username: "YLWdXJb1",
				Password: "k9cNzfB6",
			},
			userLoginResponse: RespUserLogin{
				Code: 8400,
				Msg:  "用户名或密码错误",
				Data: nil,
			},
			hasError: false,
		},
	}
	for _, item := range tests {
		response, err := item.userloginClient()
		if item.hasError {
			if err == nil {
				t.Errorf("FAILED, expect a error, got none, input: %v, output: %v\n", item.userLogin, item.userLoginResponse)
			}
		} else {
			if response != item.userLoginResponse {
				t.Errorf("FAILED, expect %v but got %v\n", item.userLoginResponse, response)
			}
		}
	}

}

// http client for UserLogin
func (t *TestItemType) userloginClient() (RespUserLogin, error) {
	// mock local
	if t.userLogin.Username == "roger1" && t.userLogin.Password == "12345678" {
		return RespUserLogin{
			Code: 8001,
			Msg:  "login success",
			Data: nil,
		}, nil
	} else {
		return RespUserLogin{
			Code: 8400,
			Msg:  "用户名或密码错误",
			Data: nil,
		}, nil
	}
}

func TestHello(t *testing.T) {
	require.NotEmpty(t, tt.Method)
	require.Equal(t, "bomb", tt.Method)
}
