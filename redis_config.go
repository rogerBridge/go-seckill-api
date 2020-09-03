package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// 定义redis pool
var pool = &redis.Pool{
	MaxIdle: 20000,
	IdleTimeout: 300 * time.Second,
	Dial: func() (conn redis.Conn, err error) {
		con, err := redis.Dial(networkType, address,
			//redis.DialReadTimeout(5*time.Second),
			//redis.DialWriteTimeout(5*time.Second),
			)
		if err!=nil {
			return nil, err
		}
		return con,nil
	},
}