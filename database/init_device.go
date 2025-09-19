package database

import (
	"go_backend/models"
	"log"

	"gorm.io/gorm"
)

func InitializeDefaultData(db *gorm.DB) {
	var count int64
	db.AutoMigrate(&models.Device{})
	db.Model(&models.Device{}).Count(&count)

	if count == 0 {
		defaultDevices := []models.Device{
			{Name: "笔记本电脑", Status: 2, Position: "杭州"},
			{Name: "无人机", Status: 0, Position: "上海"},
			{Name: "深度相机", Status: 2, Position: "云南"},
			{Name: "激光雷达", Status: 1, Position: "日本"},
		}

		result := db.Create(&defaultDevices)
		if result.Error != nil {
			log.Fatalf("Failed to insert default devices: %v", result.Error)
		} else {
			log.Println("Inserted default devices.")
		}
	}
}
