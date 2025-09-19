// database/database.go
package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// dsn := "root:1234@tcp(127.0.0.1:3306)/deviceCard?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:lz199911@tcp(127.0.0.1:3306)/deviceCard?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func GetDB() *gorm.DB {
	return DB
}
