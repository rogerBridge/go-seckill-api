package orders

import (
	"log"
	"redisplay/mysql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InsertOrders(orderNum string, userId string, productId int, purchaseNum int, orderDatetime time.Time, status string) error {
	_, err := mysql.Conn.Exec("insert orders (order_number, user_id, product_id, purchase_number, order_datetime, status) values (?, ?, ?, ?, ?, ?)", orderNum, userId, productId, purchaseNum, orderDatetime, status)
	if err != nil {
		log.Printf("insertOrders: %s", err)
		return err
	}
	return nil
}

// 更新订单信息, 暂定两种, 取消订单, 删除订单
func UpdateOrders(status string, orderNum string) error {
	_, err := mysql.Conn.Exec("update orders set status=? where order_number=?", status, orderNum)
	if err != nil {
		log.Printf("%s", err)
	}
	return nil
}

func DeleteOrders() {

}

func QueryOrders() {

}
