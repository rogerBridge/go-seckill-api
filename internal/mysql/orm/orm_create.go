package main

import (
	"log"
	"time"
)

// example for orm create
func create() {
	// Create
	// db.Create(&Product{Code: "D42", Price: 100})
	// db.Create(&User{Name: "leo", Age: 18, Birthday: time.Now()})
	// db.Select("Name").Create(&User{Name: "tong", Age: 20, Birthday: time.Now()})
	db.Omit().Create(&User{Name: "hua", Age: 26, Birthday: time.Now()})

	db.Omit().Create(&Product{Code: "hello", Price: 100})
	// db.Omit().Create(&Product{})

	usersList := []User{
		{Name: "leo00", Age: 20, Birthday: time.Now()},
		{Name: "leo01", Age: 21, Birthday: time.Now()},
		{Name: "leo02", Age: 22, Birthday: time.Now()},
	}
	db.Create(&usersList)
	for _, user := range usersList {
		log.Println(user.ID)
	}
}
