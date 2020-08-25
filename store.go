package main

// Product exported
type Product struct {
	productID   string
	StoreNum    int
	ProductName string
}

var p1 = Product{
	productID:   "10000",
	StoreNum:    200,
	ProductName: "wahaha",
}

var p2 = Product{
	productID:   "10001",
	StoreNum:    200,
	ProductName: "coca cola",
}

var p3 = Product{
	productID:   "10002",
	StoreNum:    500,
	ProductName: "挖掘机",
}

var pList = []Product{p1, p2, p3}