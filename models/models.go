package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"grpc-cms/conf"
	"log"
)

var DB *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	DB, err = gorm.Open(conf.DatabaseType, fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName))
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
