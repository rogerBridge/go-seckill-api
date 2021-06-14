package main

import "fmt"

func update() {
	// var ps []Product
	// db.Model(&Product{}).Find(&ps)
	// for _, p := range ps {
	// 	p.Code = "updated"
	// 	db.Save(&p)
	// }

	var p Product
	db.Model(&Product{}).Find(&p, "id=?", "1")
	p.Code = "mi"
	db.Save(&p)
	fmt.Println(p)

}
