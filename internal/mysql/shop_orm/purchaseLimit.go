package shop_orm

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type PurchaseLimit struct {
	SelfDefine
	ProductID              string `gorm:"product_id" json:"productID"`
	LimitNum               int    `gorm:"limit_num" json:"limitNum"`
	StartPurchaseTimeStamp int    `gorm:"start_purchase_timestamp" json:"startPurchaseTimestamp"`
	StopPurchaseTimeStamp  int    `gorm:"stop_purchase_timestamp" json:"stopPurchaseTimestamp"`
}

// Create PurchaseLimit
// 任何需要create和update的方法, 建议使用事务
func (p *PurchaseLimit) CreatePurchaseLimit(tx *gorm.DB) error {
	if err := tx.Create(p).Error; err != nil {
		log.Println("While create purchaseLimit, error: ", err)
		return err
	}
	return nil
}

// Query PurchaseLimit, 获取全部没有删掉的PurchaseLimit
func (p *PurchaseLimit) QueryPurchaseLimits() ([]*PurchaseLimit, error) {
	var results []*PurchaseLimit
	if err := conn.Model(&PurchaseLimit{}).Find(&results).Error; err != nil {
		log.Println("While query PurchaseLimits, error: ", err)
		return results, err
	}
	return results, nil
}

// 获取PurchaseLimit由商品的product_id
func (p *PurchaseLimit) QueryPurchaseLimit() (*PurchaseLimit, error) {
	var result *PurchaseLimit
	if err := conn.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).Find(&result).Error; err != nil {
		log.Println("While query PurchaseLimit, error: ", err)
		return result, err
	}
	return result, nil
}

// update PurchaseLimit by product_id
func (p *PurchaseLimit) UpdatePurchaseLimit(tx *gorm.DB) error {
	if err := tx.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).Updates(PurchaseLimit{LimitNum: p.LimitNum, StartPurchaseTimeStamp: p.StartPurchaseTimeStamp, StopPurchaseTimeStamp: p.StopPurchaseTimeStamp}).Error; err != nil {
		log.Println("UpdatePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// delete PurchaseLimit by product_id
func (p *PurchaseLimit) DeletePurchaseLimit(tx *gorm.DB) error {
	t := time.Now()
	if err := tx.Model(&PurchaseLimit{}).Where("product_id=?", p.ProductID).Update("deleted_at", t).Error; err != nil {
		log.Println("DeletePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// 检查PurchaseLimit.ProductID是否已经存在于table中
func (p *PurchaseLimit) IfPurchaseLimitExist() bool {
	var result PurchaseLimit
	conn.Model(&PurchaseLimit{}).Where("product_id = ?", p.ProductID).First(&result)
	log.Printf("result is %+v", result)
	if result.ProductID != "" {
		log.Println("PurchaseLimit的product_id已经存在于表格中")
		return true
	}
	return false
}
