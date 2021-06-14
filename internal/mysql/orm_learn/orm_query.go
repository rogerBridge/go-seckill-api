package main

import (
	"fmt"
)

func query() {
	var user User
	var users []User
	// db.First(&user, 10)
	// db.Take(&user)
	db.Table("users").Find(&user, "id = ?", "1")
	fmt.Println(user)
	db.Table("users").Where([]int{1, 2, 3, 4}).Find(&users)
	// db.Model(&User{}).Find(&user, "id = ?", "1")
	// db.Model(&User{}).First(&user, "10")
	// db.Model(&User{}).Find(&users)
	for _, value := range users {
		fmt.Println(value)
	}
	// db.Last(&user)
	// fmt.Println(user, users)
	// db.Find()
}
