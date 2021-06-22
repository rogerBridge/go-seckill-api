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
	Price       int    `gorm:"price" json:"price"`
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
	var results = make([]*Good, 128)
	// err := conn.Raw("SELECT * FROM goods WHERE deleted_at is NULL order by id").Find(&results).Error
	if err := conn.Model(&Good{}).Order("id").Find(&results).Error; err != nil {
		return results, err
	}
	// if err != nil {
	// 	log.Println("While queryGoods, err: ", err)
	// 	return results, err
	// }
	return results, nil
}

// 已知good的product_id和product_name, 更新good的inventory和product_name
func (g *Good) UpdateGood(tx *gorm.DB) error {
	if err := tx.Model(&Good{}).Where("product_id=? AND product_name=?", g.ProductID, g.ProductName).Updates(Good{ProductName: g.ProductName, Inventory: g.Inventory}).Error; err != nil {
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
