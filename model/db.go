package model

import (
	"fmt"
	"office365/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //gorm mysql
)

// Db gorm connection
var Db *gorm.DB

// Setup DbConnection
func Setup() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.DatabaseConfig.User,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port,
		config.DatabaseConfig.Database,
		config.DatabaseConfig.Charset)

	Db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	defer Db.Close()
}
