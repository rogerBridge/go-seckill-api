package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SelfDefine struct {
	gorm.Model
	Version string `gorm:"default:v0.0.0"`
}

type Product struct {
	SelfDefine // 注意, 这里之后需要使用阿里标准table: primarykey, version, is_delete, gmtCreate, gmtUpdate,
	Code       string
	Price      uint // 单位: 分
}

type User struct {
	SelfDefine
	Name     string
	Age      int
	Birthday time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("/home/leo/Source/sqlite3/test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect sqlite3")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalln("While db.AutoMigrate, ", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("While db.AutoMigrate, ", err)
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})
	// db.Create(&User{Name: "leo", Age: 18, Birthday: time.Now()})
	db.Select("Name").Create(&User{Name: "tong", Age: 20, Birthday: time.Now()})
	db.Omit().Create(&User{Name: "hua", Age: 26, Birthday: time.Now()})

	users := []User{
		{Name: "leo00", Age: 20, Birthday: time.Now()},
		{Name: "leo01", Age: 21, Birthday: time.Now()},
		{Name: "leo02", Age: 22, Birthday: time.Now()},
	}
	db.Create(&users)
	for _, user := range users {
		log.Println(user.ID)
	}
	// Read
	var product Product
	db.First(&product, "code = ?", "D42") // find product with code D42
	// db.Select(&product)

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)
}
