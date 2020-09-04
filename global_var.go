package main

var (
	//networkType = "tcp"
	//address = "127.0.0.1:6379"
	//limitNum = 2 // 限制每个用户可以购买的相同productId的商品的数量, limitNum 必须小于等于商品实际的库存
	limitNumMap = map[string]int{
		"10000": 3,
		"10001": 1,
		"10002": 1,
	}
)