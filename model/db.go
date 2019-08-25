package model

import (
	"office365/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //gorm sqlite drive
)

// DB gorm connection
var DB *gorm.DB

// Setup DbConnection
func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", "./sqlite.db")
	if err != nil {
		panic(err.Error())
	}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(20)
	DB.DB().SetMaxOpenConns(1000)
	DB.LogMode(config.RuntimeConfig.Debug)

	// migration
	DB.AutoMigrate(&Account{})
}

// CloseDB Connection
func CloseDB() {
	defer DB.Close()
}
