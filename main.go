package main

import (
	"fmt"
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
	fmt.Println("Office API server running at port 8080")

	r := routes.InitRoutes()

	r.Run()
}
