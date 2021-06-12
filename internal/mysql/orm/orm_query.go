package main

import "fmt"

func query() {
	var user User
	// var users []User
	db.First(&user, 10)

	fmt.Println(user)
	// db.Find()
}
