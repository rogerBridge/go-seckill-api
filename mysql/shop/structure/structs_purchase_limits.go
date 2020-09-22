package structure

import "time"

type PurchaseLimits struct {
	ProductId             int
	LimitNum              int
	StartPurchaseDatetime time.Time
	EndPurchaseDatetime   time.Time
}
