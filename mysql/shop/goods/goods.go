package goods

import (
	_ "github.com/go-sql-driver/mysql"
	"go_redis/mysql"
	"go_redis/mysql/shop/structure"
	"log"
)

// insert 一些数据
func InsertGoods(productId int, productName string, inventory int) error {
	// 一般情况下, Exec方法执行不需要返回值
	_, err := mysql.Conn.Exec("insert goods (product_id, product_name, inventory) values (?, ?, ?)", productId, productName, inventory)
	if err != nil {
		log.Println(err)
		return err
	}
	//rowsAffect, err := result.RowsAffected()
	//if err!=nil {
	//	log.Printf("%v", err)
	//}
	//if rowsAffect==0{
	//	return errors.New("nothing insert")
	//}
	return nil
}

// 根据商品的product_name删除商品
func DeleteGoods(productId int) error {
	_, err := mysql.Conn.Exec("update goods set is_delete=1 where product_id=?", productId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoods(productId int, productName string, inventory int) error {
	_, err := mysql.Conn.Exec("update goods set product_name=?, inventory=? where product_id=?", productName, inventory, productId)
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(result.RowsAffected())
	return nil
}

func QueryGoods() ([]*structure.Goods, error) {
	goodsList := make([]*structure.Goods, 0, 10)
	rows, err := mysql.Conn.Query("select product_id, product_name, inventory from goods where is_delete is false")
	if err != nil {
		log.Println(err)
		return goodsList, err
	}
	for rows.Next() {
		r := new(structure.Goods)
		err = rows.Scan(&r.ProductId, &r.ProductName, &r.Inventory)
		if err != nil {
			log.Println(err)
			return goodsList, err
		}
		goodsList = append(goodsList, r)
	}
	return goodsList, nil
}
