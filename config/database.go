package config

import (
	"fmt"
	"log"
	"tcm-server-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "oasis:fs@tcp(127.0.0.1:3306)/tcmDB?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// 自动迁移模型
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
	}

	fmt.Println("Database connected!")
}
