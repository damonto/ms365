package main

import (
	"office365/config"
	"office365/model"
	"office365/routes"

	"github.com/gin-gonic/gin"
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
