package main

var (
	networkType = "tcp"
	address = "172.17.0.6:6379"
	//limitNum = 2 // 限制每个用户可以购买的相同productId的商品的数量, limitNum 必须小于等于商品实际的库存
	limitNumMap = map[string]int{
		"10000": 2,
		"10001": 1,
		"10002": 1,
	}
)