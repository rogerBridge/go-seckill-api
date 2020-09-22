package structure

type Jwt struct {
	Username string `json:"username"`
	Jwt      string `json:"token"` // JWT token
}
