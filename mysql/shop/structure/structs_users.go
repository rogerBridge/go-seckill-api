package structure

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogout struct {
	Message string `json:"msg"`
}

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}
