package shop_orm

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type PurchaseLimit struct {
	SelfDefine
	ProductName            string `gorm:"column:product_name" json:"productName"`
	LimitNum               int    `gorm:"column:limit_num" json:"limitNum"`
	StartPurchaseTimeStamp int    `gorm:"column:start_purchase_timestamp" json:"startPurchaseTimestamp"`
	StopPurchaseTimeStamp  int    `gorm:"column:stop_purchase_timestamp" json:"stopPurchaseTimestamp"`
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
func (p *PurchaseLimit) QueryPurchaseLimits() ([]*PurchaseLimit, error) {
	var results []*PurchaseLimit
	if err := conn.Model(&PurchaseLimit{}).Find(&results).Error; err != nil {
		log.Println("While query PurchaseLimits, error: ", err)
		return results, err
	}
	return results, nil
}

// 获取PurchaseLimit由商品的product_name
func (p *PurchaseLimit) QueryPurchaseLimit() (*PurchaseLimit, error) {
	if !p.IfPurchaseLimitExist() {
		return p, fmt.Errorf("查找的商品不存在")
	}
	var result *PurchaseLimit
	if err := conn.Model(&PurchaseLimit{}).Where("product_name=?", p.ProductName).Find(&result).Error; err != nil {
		log.Println("While query PurchaseLimit, error: ", err)
		return result, err
	}
	return result, nil
}

// update PurchaseLimit by product_id
func (p *PurchaseLimit) UpdatePurchaseLimit(tx *gorm.DB) error {
	if err := p.CheckPurchaseLimitParams(); err != nil {
		return err
	}
	if !p.IfPurchaseLimitExist() {
		return fmt.Errorf("更新的商品不存在")
	}
	if err := tx.Model(&PurchaseLimit{}).Where("product_name=?", p.ProductName).Updates(PurchaseLimit{LimitNum: p.LimitNum, StartPurchaseTimeStamp: p.StartPurchaseTimeStamp, StopPurchaseTimeStamp: p.StopPurchaseTimeStamp}).Error; err != nil {
		log.Println("UpdatePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// delete PurchaseLimit by product_id
func (p *PurchaseLimit) DeletePurchaseLimit(tx *gorm.DB) error {
	if !p.IfPurchaseLimitExist() {
		return fmt.Errorf("删除时发现商品不存在")
	}
	t := time.Now()
	if err := tx.Model(&PurchaseLimit{}).Where("product_name=?", p.ProductName).Update("deleted_at", t).Error; err != nil {
		logger.Println("DeletePurchaseLimit error: ", err)
		return err
	}
	return nil
}

// 检查PurchaseLimit.ProductID是否已经存在于table中
func (p *PurchaseLimit) IfPurchaseLimitExist() bool {
	var result PurchaseLimit
	if err := conn.Model(&PurchaseLimit{}).Where("product_name=?", p.ProductName).First(&result).Error; err == nil {
		if result.ProductName == p.ProductName {
			return true
		}
	}
	return false
}

func (p *PurchaseLimit) CheckPurchaseLimitParams() error {
	if p.LimitNum == 0 || p.ProductName == "" {
		return fmt.Errorf("缺少参数")
	}
	return nil
}
