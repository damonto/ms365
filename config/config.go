package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// App 配置项
type App struct {
	Domain       string
	AuthorizeURL string
	TokenURL     string
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
}

// AppConfig 系统配置文件
var AppConfig = new(App)

// Runtime 配置项
type Runtime struct {
	Debug bool
}

// RuntimeConfig 运行时配置文件
var RuntimeConfig = new(Runtime)

// Database 数据库配置
type Database struct {
	Host     string
	Port     int32
	User     string
	Password string
	Database string
	Charset  string
}

// DatabaseConfig 数据库配置
var DatabaseConfig = new(Database)

// API 调用凭据
type API struct {
	AccessKey    string
	AccessSecret string
}

// APIConfig api 凭据
var APIConfig = new(API)

// Setup Config
func Setup() {
	wd, _ := os.Getwd()
	cfg, err := ini.Load(fmt.Sprintf("%v/app.ini", wd))
	if err != nil {
		fmt.Printf("fail to parse config file ./app.ini error: %v", err)
		os.Exit(1)
	}

	cfg.Section("App").MapTo(AppConfig)
	cfg.Section("Runtime").MapTo(RuntimeConfig)
	cfg.Section("Database").MapTo(DatabaseConfig)
	cfg.Section("Api").MapTo(APIConfig)
}
