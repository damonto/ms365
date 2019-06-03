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
	ClientID     string
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
}
