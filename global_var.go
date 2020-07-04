package main

var (
	orderNumLength = 10 // 839亿亿
	networkType = "tcp"
	address = "localhost:6379"
	passwd = "hello"
	//limitNum = 2 // 限制每个用户可以购买的相同productId的商品的数量, limitNum 必须小于等于商品实际的库存
	limitNumMap = map[string]int{
		"10000": 2,
		"10001": 1,
		"10002": 1,
	}
)