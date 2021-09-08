package shop_orm

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	SelfDefine
	OrderNumber string `gorm:"" json:"orderNumber"`
	Username    string `gorm:"" json:"username"`
	ProductID   int    `gorm:"" json:"productID"`
	PurchaseNum int    `gorm:"" json:"purchaseNum"`
	Status      string `gorm:"" json:"status"`
	Price       int    `gorm:"" json:"price"`
}

// check params in http request maybe is best :)
func (o *Order) CreateOrder(tx *gorm.DB) error {
	if !o.CheckOrderParams() {
		return fmt.Errorf("Order参数检查没有通过")
	}
	if err := tx.Model(&Order{}).Select("OrderNumber", "Username", "ProductID", "PurchaseNum", "Status", "Price").Create(o).Error; err != nil {
		return err
	}
	return nil
}

// find all orders
// just for admin
func (o *Order) QueryOrders() ([]*Order, error) {
	orders := make([]*Order, 128)
	if err := conn.Model(&Order{}).Find(&orders).Error; err != nil {
		return orders, err
	}
	return orders, nil
}

// find specific order, by username
func (o *Order) QueryOrderByUsername(username string) ([]*Order, error) {
	orders := make([]*Order, 128)
	if err := conn.Model(&Order{}).Where("username = ?", username).Find(&orders).Error; err != nil {
		return orders, err
	}
	return orders, nil
}

func (o *Order) UpdateOrderStatus(tx *gorm.DB) error {
	if err := tx.Model(&Order{}).Where("order_number = ?", o.OrderNumber).Update("status", o.Status).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) UpdateOrder(tx *gorm.DB) error {
	if !o.CheckOrderParams() {
		return fmt.Errorf("Order参数检查没有通过")
	}
	if err := tx.Model(&Order{}).Where("order_number=?", o.OrderNumber).Updates(Order{
		OrderNumber: o.OrderNumber,
		Username:    o.Username,
		ProductID:   o.ProductID,
		PurchaseNum: o.PurchaseNum,
		Status:      o.Status,
		Price:       o.Price,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) DeleteOrder(tx *gorm.DB) error {
	if err := tx.Model(&Order{}).Where("order_number = ?", o.OrderNumber).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) CheckOrderParams() bool {
	if o.OrderNumber == "" || o.Username == "" || o.ProductID == 0 || o.PurchaseNum <= 0 || o.Status == "" {
		return false
	}
	return true
}
