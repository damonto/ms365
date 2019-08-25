package main

import (
	"office365/config"
	"office365/model"
	"office365/routes"
)

func init() {
	config.Setup()
	model.Setup()

	if !config.RuntimeConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	r := routes.InitRoutes()

	r.Run()
}
