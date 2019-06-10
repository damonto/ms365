package main

import (
	"office365/config"
	"office365/model"
	"office365/routes"
)

func init() {
	config.Setup()
	model.Setup()
}

func main() {
	r := routes.InitRoutes()

	r.Run()
}
