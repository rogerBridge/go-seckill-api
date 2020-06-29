package main

import "github.com/gomodule/redigo/redis"

// 定义redis
var pool = &redis.Pool{
	Dial: func() (conn redis.Conn, err error) {
		return redis.Dial("tcp", "localhost:6379", redis.DialPassword("hello"))
	},
}