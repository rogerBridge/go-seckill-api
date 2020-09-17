package structure

type Goods struct {
	ProductId   int
	ProductName string
	Inventory   int
}

type GoodDelete struct {
	ProductId int
}
