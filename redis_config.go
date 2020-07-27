package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// 定义redis pool
var pool = &redis.Pool{
	MaxIdle: 10000,
	IdleTimeout: 300 * time.Second,
	Dial: func() (conn redis.Conn, err error) {
		return redis.Dial(networkType, address)
	},
}