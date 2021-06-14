package shop_orm

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Good struct {
	SelfDefine
	ProductID   string `gorm:"product_id" json:"productID"`
	ProductName string `gorm:"product_name" json:"productName"`
	Inventory   int    `gorm:"inventory" json:"inventory"`
}

// 添加记录到goods table
func (g *Good) CreateGood(tx *gorm.DB) error {
	if err := tx.Create(g).Error; err != nil {
		return err
	}
	return nil
}

// 查找记录 goods table
// 把所有deleted_at=null的实例显示出来
func (g *Good) QueryGoods() ([]*Good, error) {
	var result = make([]*Good, 128)
	err := conn.Raw("SELECT * FROM goods WHERE deleted_at is NULL order by id").Find(&result).Error
	if err != nil {
		log.Println("While queryGoods, err: ", err)
		return result, err
	}
	return result, err
}

// 已知good的product_id和product_name, 更新good的inventory
func (g *Good) UpdateGoodInventory(tx *gorm.DB) error {
	if err := tx.Model(&Good{}).Where("product_id=? AND product_name=?", g.ProductID, g.ProductName).Update("inventory", g.Inventory).Error; err != nil {
		return err
	}
	return nil
}

// 根据product_id AND product_name确定唯一的一个商品, 然后删除它
func (g *Good) DeleteGood(tx *gorm.DB) error {
	if err := tx.Model(&Good{}).Where("product_id=? AND product_name=?", g.ProductID, g.ProductName).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

// 查找某个值为productId的商品是否存在
// productID 和 productName 都必须唯一
func (g *Good) IfGoodExist() bool {
	var result Good
	conn.Model(&Good{}).Where("product_id=? AND product_name=?", g.ProductID, g.ProductName).First(&result)
	if result.ProductName != "" || result.ProductID != "" {
		log.Println("查找的商品存在")
		return true
	}
	return false
}
