package shop_orm

type Order struct {
	SelfDefine
	OrderNumber string `gorm:"order_number" json:"orderNumber"`
	Username    string `gorm:"username" json:"username"`
	ProductID   string `gorm:"product_id" json:"productID"`
	PurchaseNum int    `gorm:"purchase_num" json:"purchaseNum"`
	Status      string `gorm:"status" json:"status"`
}
