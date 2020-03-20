package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
)

var DB *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	DB, err = gorm.Open(
		viper.GetString("database.driver"),
		fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			viper.Get("database.user"),
			viper.Get("database.password"),
			viper.Get("database.host"),
			viper.Get("database.name"),
		),
	)
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	DB.SingularTable(true) // 禁用默认表名的复数形式
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer DB.Close()
}
