package goods

import (
	"database/sql"
	"go-seckill/internal/mysql"
	"go-seckill/internal/mysql/shop/structure"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// insert 一些数据
func InsertGoods(tx *sql.Tx, productId int, productName string, inventory int) error {
	// 一般情况下, 所有的sql操作都应该开启事务
	// 一般情况下, Exec方法执行不需要返回值
	_, err := tx.Exec("insert goods (product_id, product_name, inventory) values (?, ?, ?)", productId, productName, inventory)
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
func DeleteGoods(tx *sql.Tx, productId int) error {
	_, err := tx.Exec("update goods set is_delete=1 where product_id=?", productId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoods(tx *sql.Tx, productId int, productName string, inventory int) error {
	_, err := tx.Exec("update goods set product_name=?, inventory=? where product_id=?", productName, inventory, productId)
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(result.RowsAffected())
	return nil
}

// query goods from shop.goods table
// 从shop.goods中获取商品清单
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

// 查找某个值为productId的商品是否存在
func IsExist(productId int) (int, error) {
	row := mysql.Conn.QueryRow("select exists(select * from goods where product_id=? and is_delete=0)", productId)
	var isExist int
	err := row.Scan(&isExist)
	if err != nil {
		log.Printf("row.Scan error: %+v\n", err)
		return 0, err
	}
	if isExist == 1 {
		return 1, nil
	}
	return 0, nil
}
