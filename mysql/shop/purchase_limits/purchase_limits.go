package purchase_limits

import (
	"redisplay/mysql"
	"redisplay/mysql/shop/structure"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InsertPurchaseLimits(productID int, limitNum int, startPurchaseTime time.Time, endPurchaseTime time.Time) error {
	_, err := mysql.Conn.Exec("insert purchase_limits (product_id, limit_num, start_purchase_time, end_purchase_time) values (?, ?, ?, ?)", productID, limitNum, startPurchaseTime, endPurchaseTime)
	if err != nil {
		return err
	}
	return nil
}

func DeletePurchaseLimits(productID int) error {
	_, err := mysql.Conn.Exec("update purchase_limits set is_delete=1 where product_id=?", productID)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePurchaseLimits(productID int, limitNum int, startPurchaseTime time.Time, endPurchaseTime time.Time) error {
	_, err := mysql.Conn.Exec("update purchase_limits set limit_num=?, start_purchase_time=?, end_purchase_time=? where product_id=?", limitNum, startPurchaseTime, endPurchaseTime, productID)
	if err != nil {
		return err
	}
	return nil
}

func QueryPurchaseLimits() ([]*structure.PurchaseLimits, error) {
	limitsList := make([]*structure.PurchaseLimits, 0, 10)
	rows, err := mysql.Conn.Query("select product_id, limit_num, start_purchase_time, end_purchase_time from purchase_limits where is_delete is false")
	if err != nil {
		return limitsList, err
	}
	for rows.Next() {
		p := new(structure.PurchaseLimits)
		err := rows.Scan(&p.ProductId, &p.LimitNum, &p.StartPurchaseDatetime, &p.EndPurchaseDatetime)
		if err != nil {
			return limitsList, err
		}
		limitsList = append(limitsList, p)
	}
	return limitsList, nil
}
