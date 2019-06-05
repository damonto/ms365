package model

import (
	"fmt"
	"office365/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //gorm mysql
)

// DB gorm connection
var DB *gorm.DB

// Setup DbConnection
func Setup() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.DatabaseConfig.User,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port,
		config.DatabaseConfig.Database,
		config.DatabaseConfig.Charset)

	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(20)
	DB.DB().SetMaxOpenConns(1000)
	DB.LogMode(config.RuntimeConfig.Debug)
}

// CloseDB Connection
func CloseDB() {
	defer DB.Close()
}
