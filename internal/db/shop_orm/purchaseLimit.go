package shop_orm

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// 想了好久, 还是用product_id吧, 就是redis那边的转换可能需要改不少代码
// 如果我直接在goodList这个接口中, 把ID的输出改为字符串不就好了?
type PurchaseLimit struct {
	SelfDefine
	ProductID              int `gorm:"column:product_id" json:"productID"`
	LimitNum               int `gorm:"column:limit_num" json:"limitNum"`
	StartPurchaseTimeStamp int `gorm:"column:start_purchase_timestamp" json:"startPurchaseTimestamp"`
	StopPurchaseTimeStamp  int `gorm:"column:stop_purchase_timestamp" json:"stopPurchaseTimestamp"`
}

// Create PurchaseLimit
// 任何需要create和update的方法, 建议使用事务
func (p *PurchaseLimit) CreatePurchaseLimit(tx *gorm.DB) error {
	if err := p.CheckPurchaseLimitParams(); err != nil {
		return err
	}
	if p.IfPurchaseLimitExist() {
		return fmt.Errorf("purchaseLimit existed")
	}
	if err := tx.Create(p).Error; err != nil {
		log.Println("While create purchaseLimit, error: ", err)
		return err
	}
	return nil
}

// Query PurchaseLimit, 获取全部的没有删掉的PurchaseLimit
func (p *PurchaseLimit) QueryPurchaseLimits() []*PurchaseLimit {
	var results []*PurchaseLimit
	conn.Model(&PurchaseLimit{}).Find(&results)
	return results
}

// 获取PurchaseLimit由商品的product_id
func (p *PurchaseLimit) QueryPurchaseLimitByProductID() *PurchaseLimit {
	// if !p.IfPurchaseLimitExist() {
	// 	return p, fmt.Errorf("查找的商品不存在")
	// }
	var result *PurchaseLimit
	conn.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).Find(&result)
	if p.ProductID != result.ProductID {
		return nil
	}
	return result

}

// update PurchaseLimit by product_id
func (p *PurchaseLimit) UpdatePurchaseLimit(tx *gorm.DB) error {
	if err := p.CheckPurchaseLimitParams(); err != nil {
		return err
	}
	if !p.IfPurchaseLimitExist() {
		return fmt.Errorf("更新的商品不存在")
	}
	if err := tx.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).Updates(PurchaseLimit{LimitNum: p.LimitNum, StartPurchaseTimeStamp: p.StartPurchaseTimeStamp, StopPurchaseTimeStamp: p.StopPurchaseTimeStamp}).Error; err != nil {
		log.Println("UpdatePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// delete PurchaseLimit by product_id
func (p *PurchaseLimit) DeletePurchaseLimit(tx *gorm.DB) error {
	if p.ProductID == 0 {
		return fmt.Errorf("缺少参数")
	}
	if !p.IfPurchaseLimitExist() {
		return fmt.Errorf("删除时发现商品不存在")
	}
	if err := tx.Model(&PurchaseLimit{}).Delete(p).Error; err != nil {
		logger.Println("DeletePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// 检查PurchaseLimit.ProductID是否已经存在于table中
func (p *PurchaseLimit) IfPurchaseLimitExist() bool {
	var result PurchaseLimit
	conn.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).First(&result)
	if p.ProductID == result.ProductID {
		return true
	} else {
		return false
	}
}

func (p *PurchaseLimit) CheckPurchaseLimitParams() error {
	if p.LimitNum == 0 || p.ProductID == 0 {
		return fmt.Errorf("缺少参数")
	}
	return nil
}
