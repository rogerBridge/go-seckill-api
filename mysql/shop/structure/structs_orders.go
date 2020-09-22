package structure

import "time"

type Orders struct {
	OrderNum      string
	UserId        string
	ProductId     int
	PurchaseNum   int
	OrderDatetime time.Time
	Status        string
}
