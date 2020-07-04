package main

type Product struct {
	ProductId   string
	StoreNum    int
	ProductName string
}

var p1 = Product{
	ProductId:   "10000",
	StoreNum:    200,
	ProductName: "wahaha",
}

var p2 = Product{
	ProductId:   "10001",
	StoreNum:    200,
	ProductName: "coca cola",
}

var p3 = Product{
	ProductId:   "10002",
	StoreNum:    100,
	ProductName: "挖掘机",
}

var pList = []Product{p1, p2, p3}