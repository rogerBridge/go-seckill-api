package main

func delete() {
	db.Table("users").Where("name=?", "hua").Delete(1)
}
