package structure

import "time"

type Jwt struct {
	Username     string    `json:"username"`
	Jwt          string    `json:"token"` // JWT token
	GenerateTime time.Time `json:"generateTime"`
	ExpireTime   time.Time `json:"expireTime"`
}
